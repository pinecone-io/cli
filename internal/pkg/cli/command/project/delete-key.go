package project

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type DeleteApiKeyOptions struct {
	name string
	yes  bool
}

func NewDeleteKeyCmd() *cobra.Command {
	options := DeleteApiKeyOptions{}

	cmd := &cobra.Command{
		Use:     "delete-key",
		Short:   "delete an API key in a project",
		GroupID: help.GROUP_PROJECTS_API_KEYS.ID,
		Example: help.Examples([]string{
			"pinecone target -o \"my-org\" -p \"my-project\"",
			"pinecone delete-key -n \"my-key\"",
		}),
		Run: func(cmd *cobra.Command, args []string) {
			projId, err := getTargetProjectId()
			if err != nil {
				msg.FailMsg("No target project set. Use %s to set the target project.", style.Code("pinecone target -o <org> -p <project>"))
				exit.ErrorMsg("No project context set")
			}

			if options.name == "" {
				msg.FailMsg("Name of the key is required. Pass it with the %s flag", style.Code("--name"))
				exit.ErrorMsg("Name of the key is required")
			}

			// Verify key exists before trying to delete it.
			// This lets us give a more helpful error message than just
			// attempting to delete non-existant key and getting 500 error.
			existingKeys, err := dashboard.GetApiKeysById(projId)
			if err != nil {
				msg.FailMsg("Failed to list keys: %s", err)
				exit.Error(err)
			}
			var keyToDelete dashboard.Key
			var keyExists bool = false
			for _, key := range existingKeys.Keys {
				if key.UserLabel == options.name {
					keyToDelete = key
					keyExists = true
				}
			}
			if !keyExists {
				msg.FailMsg("Key with name %s does not exist", style.Emphasis(options.name))
				msg.HintMsg("See existing keys with %s", style.Code("pinecone project list-keys"))
				exit.ErrorMsg(pcio.Sprintf("Key with name %s does not exist", style.Emphasis(options.name)))
			}

			if !options.yes {
				confirmDeleteApiKey(options.name)
			}

			_, err = dashboard.DeleteApiKey(projId, keyToDelete)
			if err != nil {
				msg.FailMsg("Failed to delete key: %s", err)
				exit.Error(err)
			}
			msg.SuccessMsg("API key %s deleted", style.Emphasis(options.name))
		},
	}

	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the key to create")
	cmd.Flags().BoolVar(&options.yes, "yes", false, "skip confirmation prompt")
	return cmd
}

func confirmDeleteApiKey(apiKeyName string) {
	msg.WarnMsg("This operation will delete API Key %s from project %s.", style.Emphasis(apiKeyName), style.Emphasis(state.TargetProj.Get().Name))
	msg.WarnMsg("Any integrations you have that auth with this API Key will stop working.")
	msg.WarnMsg("This action cannot be undone.")

	// Prompt the user
	fmt.Print("Do you want to continue? (y/N): ")

	// Read the user's input
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Trim any whitespace from the input and convert to lowercase
	input = strings.TrimSpace(strings.ToLower(input))

	// Check if the user entered "y" or "yes"
	if input == "y" || input == "yes" {
		msg.InfoMsg("You chose to continue delete.")
	} else {
		msg.InfoMsg("Operation canceled.")
		exit.Success()
	}
}
