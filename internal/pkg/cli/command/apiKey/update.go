package apiKey

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type UpdateAPIKeyOptions struct {
	apiKeyId string

	name  string
	roles []string

	json bool
}

func NewUpdateAPIKeyCmd() *cobra.Command {
	options := UpdateAPIKeyOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing API key by ID with the specified configuration",
		Example: help.Examples(`
			pc api-key update --id "api-key-id" --name "updated-name" --roles "ProjectEditor"
		`),
		GroupID: help.GROUP_API_KEYS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			// Only set non-empty values
			updateParams := &pinecone.UpdateAPIKeyParams{}
			if options.name != "" {
				updateParams.Name = &options.name
			}
			if options.roles != nil {
				updateParams.Roles = &options.roles
			}

			apiKey, err := ac.APIKey.Update(cmd.Context(), options.apiKeyId, updateParams)
			if err != nil {
				msg.FailMsg("Failed to update API key %s: %s\n", options.apiKeyId, err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(apiKey)
				pcio.Println(json)
				return
			}

			msg.SuccessMsg("API key %s updated successfully.", apiKey.Id)
			presenters.PrintDescribeAPIKeyTable(apiKey)
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.apiKeyId, "id", "i", "", "id of the API key to update")
	_ = cmd.MarkFlagRequired("id")

	// optional flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "The new name for the API key")
	cmd.Flags().StringSliceVarP(&options.roles, "roles", "r", []string{}, "The new roles for the API key")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}
