package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

type pruneLocalKeysCmdOptions struct {
	origin           string // "cli", "user", "all"
	projectID        string // optional filter, if not provided all projects will be pruned
	dryRun           bool   // preview keys that will be deleted
	skipConfirmation bool   // skip confirmation prompt
	json             bool
}

var (
	pruneHelp = help.Long(`
		Delete project API keys that the CLI is managing in local state.

		This operation removes managed keys from local storage and deletes them
		from Pinecone servers. Any integrations that authenticate with these
		keys outside of the CLI will immediately stop working.

		By default, this command deletes all API keys that the CLI is managing, whether
		they were created by the CLI or the user. You can filter the operation by 
		project ID or key origin. Options for origin are: 'cli', 'user', or 'all' (default).

		See: https://docs.pinecone.io/reference/cli/authentication
	`)

	pruneExample = help.Examples(`
		# Prune all locally managed keys that the CLI has created
		pc auth local-keys prune --origin cli --skip-confirmation

		# Prune all locally managed keys that the user has created and stored
		pc auth local-keys prune --origin user --skip-confirmation

		# Show a dry run plan of pruning all keys (origin defaults to "all")
		# and then apply the changes
		pc auth local-keys prune --dry-run --skip-confirmation
		pc auth local-keys prune --skip-confirmation
	`)
)

func NewPruneLocalKeysCmd() *cobra.Command {
	options := pruneLocalKeysCmdOptions{}

	cmd := &cobra.Command{
		Use:     "prune",
		Short:   "Delete project API keys that the CLI is managing in local storage",
		Long:    pruneHelp,
		Example: pruneExample,
		Run: func(cmd *cobra.Command, args []string) {
			runPruneLocalKeys(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.origin, "origin", "o", "all", "Filter deletions by key origin: 'cli', 'user', 'all'")
	cmd.Flags().StringVar(&options.projectID, "id", "", "Only prune keys for a specific project")
	cmd.Flags().BoolVar(&options.dryRun, "dry-run", false, "Preview keys that will be deleted without applying changes")
	cmd.Flags().BoolVar(&options.skipConfirmation, "skip-confirmation", false, "Skip confirmation prompt")
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}

func runPruneLocalKeys(ctx context.Context, options pruneLocalKeysCmdOptions) {
	ac := sdk.NewPineconeAdminClient()
	managedKeys := secrets.GetManagedProjectKeys()

	// Filter to projectId if provided
	if options.projectID != "" {
		if mk, ok := managedKeys[options.projectID]; ok {
			managedKeys = map[string]secrets.ManagedKey{options.projectID: mk}
		} else {
			msg.FailMsg("No managed keys found for project ID %s", style.Emphasis(options.projectID))
			exit.ErrorMsgf("no managed keys found for project ID: %s", options.projectID)
		}
	}

	// Build dry run plan
	var plan []planItem
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
			plan = append(plan, planItem{projectID: projectID, managedKey: managedKey, onServer: onServer})
		}
	}

	// If there's nothing in the plan, we can exit
	if len(plan) == 0 {
		msg.InfoMsg("No locally managed API keys to prune")
		return
	}

	// Dry run preview
	if options.dryRun {
		printDryRunPlan(plan, options)
		msg.InfoMsg("Dry run complete. Re-run with %s and %s to apply changes", style.Emphasis("--yes"), style.Emphasis("--dry-run=false"))
		return
	}

	// Confirm if we should apply pruning changes
	shouldPrune := true
	if !options.skipConfirmation {
		confirmed, err := confirmPruneKeys(plan, options)
		if err != nil {
			msg.FailMsg("Failed to confirm pruning keys: %s", err)
			exit.Error(err, "Failed to confirm pruning keys")
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
				exit.Errorf(err, "Failed to delete remote key %s", key.managedKey.Id)
				continue // If we failed to delete the remote key, move on and keep the locally stored key for now
			}
			msg.SuccessMsg("Deleted remote key %s (project %s)", style.Emphasis(key.managedKey.Id), style.Emphasis(key.projectID))
		}

		// Clean up local record
		secrets.DeleteProjectManagedKey(key.projectID)
		msg.SuccessMsg("Deleted local record for key %s (project %s)", style.Emphasis(key.managedKey.Id), style.Emphasis(key.projectID))
	}

	msg.SuccessMsg("Pruning operation complete")
}

func includeByOrigin(origin secrets.ManagedKeyOrigin, filter string) bool {
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

func confirmPruneKeys(plan []planItem, options pruneLocalKeysCmdOptions) (bool, error) {
	msg.WarnMsg("This operation will delete the following API Keys:")
	printDryRunPlan(plan, options)
	msg.WarnMsg("Any integrations you have that authenticate with these API keys will immediately stop working.")
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

func printDryRunPlan(plan []planItem, options pruneLocalKeysCmdOptions) {
	if options.json {
		json := text.IndentJSON(plan)
		pcio.Println(json)
	} else {
		for _, key := range plan {
			if key.onServer {
				msg.WarnMsg("API key %s (ID: %s) will be deleted locally and remotely from project %s (ID: %s)",
					style.Emphasis(key.managedKey.Name),
					style.Emphasis(key.managedKey.Id),
					style.Emphasis(key.managedKey.ProjectName),
					style.Emphasis(key.projectID))
			} else {
				msg.WarnMsg("Local API key %s (ID: %s) was not found in project %s (ID: %s) and will be deleted locally ",
					style.Emphasis(key.managedKey.Name),
					style.Emphasis(key.managedKey.Id),
					style.Emphasis(key.managedKey.ProjectName),
					style.Emphasis(key.projectID))
			}
		}
	}
}

type planItem struct {
	projectID  string
	managedKey secrets.ManagedKey
	onServer   bool
}
