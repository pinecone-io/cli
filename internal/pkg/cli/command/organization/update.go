package organization

import (
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

type UpdateOrganizationCmdOptions struct {
	organizationID string
	name           string

	json bool
}

func NewUpdateOrganizationCmd() *cobra.Command {
	options := UpdateOrganizationCmdOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing organization by ID, or the target organization",
		Example: help.Examples(`
			pc organization update --id "organization-id" --name "new-name"
		`),
		GroupID: help.GROUP_ORGANIZATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			orgId := options.organizationID
			var err error
			if orgId == "" {
				orgId, err = state.GetTargetOrgId()
				if err != nil {
					msg.FailMsg("No target organization set and no organization ID provided. Use %s to set the target organization. Use %s to describe an organization by ID.", style.Code("pc target -o <org>"), style.Code("pc organization describe -i <organization-id>"))
					exit.ErrorMsg("No organization ID provided, and no target organization set")
				}
			}

			// Only set non-empty values
			updateParams := &pinecone.UpdateOrganizationParams{}
			if options.name != "" {
				updateParams.Name = &options.name
			}

			org, err := ac.Organization.Update(cmd.Context(), orgId, updateParams)
			if err != nil {
				msg.FailMsg("Failed to update organization %s: %s\n", orgId, err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(org)
				pcio.Println(json)
				return
			}

			msg.SuccessMsg("Organization %s updated successfully.", org.Id)
			presenters.PrintDescribeOrganizationTable(org)
		},
	}

	// optional flags
	cmd.Flags().StringVarP(&options.organizationID, "id", "i", "", "The ID of the organization to update if not the target organization")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "The new name to use for the organization")
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}
