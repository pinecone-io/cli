package apiKey

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/interactive"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type DeleteApiKeyOptions struct {
	apiKeyId string
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
		$ pc api-key delete -i "api-key-id" -y
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

			// Check if -y flag is set
			assumeYes, _ := cmd.Flags().GetBool("assume-yes")
			if !assumeYes {
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

	return cmd
}

func confirmDeleteApiKey(apiKeyName string) {
	msg.WarnMsgMultiLine(
		pcio.Sprintf("This operation will delete API Key %s from project %s.", style.Emphasis(apiKeyName), style.Emphasis(state.TargetProj.Get().Name)),
		"Any integrations you have that auth with this API Key will stop working.",
		"This action cannot be undone.",
	)

	question := "Are you sure you want to proceed with deleting this API key?"
	if !interactive.GetConfirmation(question) {
		msg.InfoMsg("Operation canceled.")
		exit.Success()
	}
	msg.InfoMsg("You chose to continue delete.")
}
