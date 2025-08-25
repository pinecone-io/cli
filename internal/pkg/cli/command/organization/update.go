package organization

import (
	"github.com/MakeNowJust/heredoc"
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

type UpdateOrganizationCmdOptions struct {
	organizationID string
	name           string

	json bool
}

func NewUpdateOrganizationCmd() *cobra.Command {
	options := UpdateOrganizationCmdOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing organization by ID with the specified configuration",
		Example: heredoc.Doc(`
		$ pc organization update -i <organization-id> --n <new-name>
		`),
		GroupID: help.GROUP_ORGANIZATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			updateParams := &pinecone.UpdateOrganizationParams{}

			if options.name != "" {
				updateParams.Name = &options.name
			}

			org, err := ac.Organization.Update(cmd.Context(), options.organizationID, updateParams)
			if err != nil {
				msg.FailMsg("Failed to update organization %s: %s\n", options.organizationID, err)
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

	// required flags
	cmd.Flags().StringVarP(&options.organizationID, "id", "i", "", "The ID of the organization to update")
	_ = cmd.MarkFlagRequired("id")

	// optional flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "The new name to use for the organization")
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}
