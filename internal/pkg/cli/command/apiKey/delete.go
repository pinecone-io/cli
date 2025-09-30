package apiKey

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type DeleteApiKeyOptions struct {
	apiKeyId         string
	skipConfirmation bool
}

func NewDeleteKeyCmd() *cobra.Command {
	options := DeleteApiKeyOptions{}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete an API key by ID",
		GroupID: help.GROUP_API_KEYS.ID,
		Example: help.Examples(`
			pc api-key delete -i "api-key-id" 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			// Verify key exists before trying to delete it.
			// This lets us give a more helpful error message than just
			// attempting to delete non-existent key and getting 500 error.
			keyToDelete, err := ac.APIKey.Describe(cmd.Context(), options.apiKeyId)
			if err != nil {
				msg.FailMsg("Failed to describe existing API key: %s", err)
				exit.Error(err)
			}

			if !options.skipConfirmation {
				confirmDeleteApiKey(keyToDelete.Name)
			}

			err = ac.APIKey.Delete(cmd.Context(), keyToDelete.Id)
			if err != nil {
				msg.FailMsg("Failed to delete API key %s: %s", style.Emphasis(keyToDelete.Name), err)
				exit.Error(err)
			}
			msg.SuccessMsg("API key %s deleted", style.Emphasis(keyToDelete.Name))

			// Check if the key is locally stored and clean it up if so
			managedKey, ok := secrets.GetProjectManagedKey(keyToDelete.ProjectId)
			if ok && managedKey.Id == keyToDelete.Id {
				secrets.DeleteProjectManagedKey(keyToDelete.ProjectId)
				msg.SuccessMsg("Deleted local record for key %s (project %s)", style.Emphasis(keyToDelete.Id), style.Emphasis(keyToDelete.ProjectId))
			}
		},
	}

	cmd.Flags().StringVarP(&options.apiKeyId, "id", "i", "", "The ID of the API key to delete")
	_ = cmd.MarkFlagRequired("id")

	cmd.Flags().BoolVar(&options.skipConfirmation, "skip-confirmation", false, "Skip deletion confirmation prompt")
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
