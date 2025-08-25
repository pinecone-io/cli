package organization

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
	"github.com/spf13/cobra"
)

type DescribeOrganizationCmdOptions struct {
	organizationID string
	json           bool
}

func NewDescribeOrganizationCmd() *cobra.Command {
	options := DescribeOrganizationCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe an organization by ID or the target organization",
		Example: heredoc.Doc(`
		$ pc organization describe -i <organization-id>
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

			org, err := ac.Organization.Describe(cmd.Context(), orgId)
			if err != nil {
				msg.FailMsg("Failed to describe organization %s: %s\n", orgId, err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(org)
				pcio.Println(json)
			} else {
				presenters.PrintDescribeOrganizationTable(org)
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.organizationID, "id", "i", "", "The ID of the organization to describe")
	_ = cmd.MarkFlagRequired("id")

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}
