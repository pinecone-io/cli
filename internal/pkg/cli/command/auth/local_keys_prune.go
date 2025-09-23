package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type PruneLocalKeysCmdOptions struct {
	origin           string // "cli", "user", "all"
	projectId        string // optional filter, if not provided all projects will be pruned
	dryRun           bool   // preview keys that will be deleted
	skipConfirmation bool   // skip confirmation prompt
}

func NewPruneLocalKeysCmd() *cobra.Command {
	options := PruneLocalKeysCmdOptions{}

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Clean up project API keys that the CLI is managing",
		Run: func(cmd *cobra.Command, args []string) {
			runPruneLocalKeys(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.origin, "origin", "o", "all", "Filter deletions by key origin: 'cli', 'user', 'all'")
	cmd.Flags().StringVar(&options.projectId, "id", "", "Only prune keys for a specific project")
	cmd.Flags().BoolVar(&options.dryRun, "dry-run", false, "Preview keys that will be deleted without applying changes")
	cmd.Flags().BoolVar(&options.skipConfirmation, "skip-confirmation", false, "Skip confirmation prompt")

	return cmd
}

func runPruneLocalKeys(ctx context.Context, options PruneLocalKeysCmdOptions) {
	ac := sdk.NewPineconeAdminClient()
	managedKeys := secrets.GetManagedProjectKeys()

	// Filter to projectId if provided
	if options.projectId != "" {
		if mk, ok := managedKeys[options.projectId]; ok {
			managedKeys = map[string]secrets.ManagedKey{options.projectId: mk}
		} else {
			msg.FailMsg("No managed keys found for project ID %s", style.Emphasis(options.projectId))
			exit.Error(pcio.Errorf("no managed keys found for project ID: %s", options.projectId))
		}
	}

	// Build dry run plan
	for projectID, managedKey := range managedKeys {
		// Check key origin and what was passed in --origin
		if !includeByOrigin(managedKey.Origin, options.origin) {
			continue
		}

		// Fetch project keys to check for orphans
		projKeys, err := ac.APIKey.List(ctx, projectID)
		if err != nil { // If we errored on fetching the project keys, skip to the next project
			log.Error().Err(err).Msg(fmt.Sprintf("Failed to list API keys for project %s: %s", style.Emphasis(projectID), err))
			msg.FailMsg("Failed to list API keys for project %s: %s", style.Emphasis(projectID), err)
			continue
		}
		projKeysMap := createKeysMap(projKeys)
		_, onServer := projKeysMap[managedKey.Id]
		if onServer {
			plan = append(plan, planItem{projectId: projectID, managedKey: managedKey, onServer: onServer})
		}
	}

	// If there's nothing in the plan, we can exit
	if len(plan) == 0 {
		msg.InfoMsg("No locally managed API keys to prune")
		return
	}

	// Dry run preview
	if options.dryRun {
		printDryRunPlan(plan)
		msg.InfoMsg("Dry run complete. Re-run with %s and %s to apply changes", style.Emphasis("--yes"), style.Emphasis("--dry-run=false"))
		return
	}

	// Confirm if we should apply pruning changes
	shouldPrune := true
	if !options.skipConfirmation {
		confirmed, err := confirmPruneKeys(plan)
		if err != nil {
			msg.FailMsg("Failed to confirm pruning keys: %s", err)
			exit.Error(err)
		}
		shouldPrune = confirmed
	}

	if !shouldPrune {
		msg.InfoMsg("Pruning operation canceled")
		return
	}

	// Apply pruning changes
	for _, key := range plan {
		if key.onServer {
			if err := ac.APIKey.Delete(ctx, key.managedKey.Id); err != nil {
				msg.FailMsg("Failed to delete remote key %s: %v", style.Emphasis(key.managedKey.Id), err)
				exit.Error(err)
				continue // If we failed to delete the remote key, move on and keep the locally stored key for now
			}
			msg.SuccessMsg("Deleted remote key %s (project %s)", style.Emphasis(key.managedKey.Id), style.Emphasis(key.projectId))
		}

		// Clean up local record
		secrets.DeleteProjectManagedKey(key.projectId)
		msg.SuccessMsg("Deleted local record for key %s (project %s)", style.Emphasis(key.managedKey.Id), style.Emphasis(key.projectId))
	}

	msg.SuccessMsg("Pruning operation complete")
}

func includeByOrigin(origin secrets.ManagedKeyOrigin, filter string) bool {
	fmt.Printf("FILTER: %s, ORIGIN: %s\n", filter, origin)
	switch strings.ToLower(filter) {
	case "cli":
		return origin == secrets.OriginCLICreated
	case "user":
		return origin == secrets.OriginUserCreated
	case "all", "":
		return true
	default:
		return true // Unknown filter
	}
}

func createKeysMap(keys []*pinecone.APIKey) map[string]struct{} {
	keysMap := make(map[string]struct{})
	for _, key := range keys {
		keysMap[key.Id] = struct{}{}
	}
	return keysMap
}

func confirmPruneKeys(plan []planItem) (bool, error) {
	msg.WarnMsg("This operation will delete the following API Keys:")
	printDryRunPlan(plan)
	msg.WarnMsg("Any integrations you have that auth with these API Keys will stop working.")
	msg.WarnMsg("This action cannot be undone.")

	// Prompt the user
	fmt.Print("Do you want to continue? (y/N): ")

	// Read the user's input
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false, err
	}

	// Trim any whitespace from the input and convert to lowercase
	input = strings.TrimSpace(strings.ToLower(input))

	// Check if the user entered "y" or "yes"
	if input == "y" || input == "yes" {
		return true, nil
	} else {
		return false, nil
	}
}

func printDryRunPlan(plan []planItem) {
	for _, key := range plan {
		if key.onServer {
			msg.WarnMsg("Would delete remote key %s and local record (project %s)",
				style.Emphasis(key.managedKey.Id),
				style.Emphasis(key.projectId))
		} else {
			msg.WarnMsg("Would delete local record for key %s (not found on server, project %s)",
				style.Emphasis(key.managedKey.Id),
				style.Emphasis(key.projectId))
		}
	}
}

type planItem struct {
	projectId  string
	managedKey secrets.ManagedKey
	onServer   bool
}

var plan []planItem
