package organization

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

func NewOrganizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "organization <command>",
		Short:   "Manage Pinecone organizations",
		GroupID: help.GROUP_ADMIN.ID,
	}

	cmd.AddGroup(help.GROUP_ORGANIZATIONS)

	cmd.AddCommand(NewListOrganizationsCmd())
	cmd.AddCommand(NewUpdateOrganizationCmd())
	cmd.AddCommand(NewDeleteOrganizationCmd())
	cmd.AddCommand(NewDescribeOrganizationCmd())

	return cmd
}
