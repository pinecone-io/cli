package apiKey

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type CreateApiKeyOptions struct {
	projectId string
	name      string
	roles     []string
	json      bool
}

func NewCreateApiKeyCmd() *cobra.Command {
	options := CreateApiKeyOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create an API key in a project",
		GroupID: help.GROUP_API_KEYS.ID,
		Example: heredoc.Doc(`
		$ pc target -o "my-org" -p "my-project"
		$ pc api-key create -n "my-key" 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			projId := options.projectId
			var err error
			if projId == "" {
				projId, err = state.GetTargetProjectId()
				if err != nil {
					msg.FailMsg("No target project set, and no project ID provided. Use %s to set the target project. Use %s to create the key in a specific project.", style.Code("pc target -o <org> -p <project>"), style.Code("pc api-key create -i <project-id> -n <name>"))
					exit.ErrorMsg("No project ID provided, and no target project set")
				}
			}

			createParams := &pinecone.CreateAPIKeyParams{}

			// Only set non-empty values
			if options.name != "" {
				createParams.Name = options.name
			}
			if options.roles != nil {
				createParams.Roles = &options.roles
			}

			keyWithSecret, err := ac.APIKey.Create(cmd.Context(), projId, createParams)
			if err != nil {
				msg.FailMsg("Failed to create API key %s in project %s: %s", options.name, projId, err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(keyWithSecret)
				pcio.Println(json)
			} else {
				msg.SuccessMsg("API Key %s created successfully.\n", style.Emphasis(keyWithSecret.Key.Name))
				presenters.PrintDescribeAPIKeyWithSecretTable(keyWithSecret)
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the key to create")
	_ = cmd.MarkFlagRequired("name")

	// optional flags
	cmd.Flags().StringVarP(&options.projectId, "id", "i", "", "ID of the project to create the key for if not the target project")
	cmd.Flags().StringSliceVar(&options.roles, "roles", []string{}, "roles to assign to the key. The default is 'ProjectEditor'")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
