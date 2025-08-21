package apiKey

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type DeleteApiKeyOptions struct {
	apiKeyId string
	yes      bool
}

func NewDeleteKeyCmd() *cobra.Command {
	options := DeleteApiKeyOptions{}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete an API key by ID",
		GroupID: help.GROUP_API_KEYS.ID,
		Example: heredoc.Doc(`
		$ pc target -o "my-org" -p "my-project"
		$ pc api-key delete -i "api-key-id" 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			projId, err := state.GetTargetProjectId()
			if err != nil {
				msg.FailMsg("No target project set. Use %s to set the target project.", style.Code("pc target -o <org> -p <project>"))
				exit.ErrorMsg("No project context set")
			}

			// Verify key exists before trying to delete it.
			// This lets us give a more helpful error message than just
			// attempting to delete non-existent key and getting 500 error.
			existingKeys, err := ac.APIKey.List(cmd.Context(), projId)
			if err != nil {
				msg.FailMsg("Failed to list keys: %s", err)
				exit.Error(err)
			}
			var keyToDelete *pinecone.APIKey
			var keyExists bool = false
			for _, key := range existingKeys {
				if key.Id == options.apiKeyId {
					keyToDelete = key
					keyExists = true
				}
			}
			if !keyExists {
				msg.FailMsg("Key with ID %s does not exist", style.Emphasis(options.apiKeyId))
				msg.HintMsg("See existing keys with %s", style.Code(pcio.Sprintf("pc api-key list")))
				exit.ErrorMsg(pcio.Sprintf("Key with ID %s does not exist", style.Emphasis(options.apiKeyId)))
			}

			if !options.yes {
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

	cmd.Flags().StringVarP(&options.apiKeyId, "id", "i", "", "the ID of the API key to delete")
	_ = cmd.MarkFlagRequired("id")

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
