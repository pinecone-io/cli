package apiKey

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
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
	store     bool
	roles     []string
	json      bool
}

func NewCreateApiKeyCmd() *cobra.Command {
	options := CreateApiKeyOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create an API key for a specific project by ID or the target project",
		GroupID: help.GROUP_API_KEYS.ID,
		Example: help.Examples(`
		    # Create a new API key for the target project
			pc target --org "org-name" --project "project-name"
			pc api-key create --name "key-name" 

			# Create a new API key for a specific project
			pc api-key create --id "project-id" --name "key-name"
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

			targetOrgId, err := state.GetTargetOrgId()
			if err != nil {
				msg.FailMsg("Failed to get target organization ID: %s", err)
				exit.Error(err)
			}

			// Only set non-empty values
			createParams := &pinecone.CreateAPIKeyParams{}
			if options.name != "" {
				createParams.Name = options.name
			}
			if len(options.roles) > 0 {
				createParams.Roles = &options.roles
			} else {
				// Default to 'ProjectEditor' role if no roles are provided
				createParams.Roles = &[]string{"ProjectEditor"}
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

			// If the user requested to store the key locally
			if options.store {
				// If a key for the project is already stored locally, delete it if its CLI managed
				managedKey, ok := secrets.GetProjectManagedKey(keyWithSecret.Key.ProjectId)
				if ok && managedKey.Origin == secrets.OriginCLICreated {
					err := ac.APIKey.Delete(cmd.Context(), managedKey.Id)
					if err != nil {
						msg.FailMsg("Failed to delete previously managed API key: %s, %+v", style.Emphasis(managedKey.Id), err)
					}
					msg.SuccessMsg("Deleted previously managed API key: %s", style.Emphasis(managedKey.Id))
				}

				// Store the new key
				msg.SuccessMsg("Storing key %s locally for future CLI operations", style.Emphasis(keyWithSecret.Key.Name))
				secrets.SetProjectManagedKey(secrets.ManagedKey{
					Name:           keyWithSecret.Key.Name,
					Id:             keyWithSecret.Key.Id,
					Value:          keyWithSecret.Value,
					ProjectId:      keyWithSecret.Key.ProjectId,
					OrganizationId: targetOrgId,
					Origin:         secrets.OriginUserCreated,
				})
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "Name of the key to create")
	_ = cmd.MarkFlagRequired("name")

	// optional flags
	cmd.Flags().StringVarP(&options.projectId, "id", "i", "", "ID of the project to create the key for if not the target project")
	cmd.Flags().BoolVar(&options.store, "store", false, "Stores the created key locally so it can be used for future CLI operations")
	cmd.Flags().StringSliceVar(&options.roles, "roles", []string{}, "Roles to assign to the key. The default is 'ProjectEditor'")
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}
