package apiKey

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/interactive"
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
		Example: heredoc.Doc(`
		$ pc target -o "my-org" -p "my-project"
		$ pc api-key delete -i "api-key-id" 
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

	question := fmt.Sprintf("Do you want to continue deleting API key '%s'?", apiKeyName)
	if !interactive.GetConfirmation(question) {
		msg.InfoMsg("Operation canceled.")
		exit.Success()
	}
	msg.InfoMsg("You chose to continue delete.")
}
