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
		Short: "Describe an organization by ID",
		Example: heredoc.Doc(`
		$ pc organization describe -i <organization-id>
		`),
		GroupID: help.GROUP_ORGANIZATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			org, err := ac.Organization.Describe(cmd.Context(), options.organizationID)
			if err != nil {
				msg.FailMsg("Failed to describe organization %s: %s\n", options.organizationID, err)
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
