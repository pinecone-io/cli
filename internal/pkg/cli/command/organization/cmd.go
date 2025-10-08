package organization

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	organizationHelp = help.Long(`
		Manage Pinecone organizations.

		An organization is a collection of projects. Organizations allow one or more users to
		control billing and project permissions for all of the projects belonging to the organization.

		When authenticating with Pinecone through the CLI, organizations are intrinsically linked to 
		user credentials. If you have authenticated using pc auth login, you will be able to target, list
		and manage any organizations linked to that account. If you have authenticated using a service account,
		or an explicit API key, you will be scoped to the organization associated with that account.

		See: https://docs.pinecone.io/guides/organizations/understanding-organizations
	`)
)

func NewOrganizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "organization <command>",
		Short:   "Manage Pinecone organizations",
		Long:    organizationHelp,
		GroupID: help.GROUP_ADMIN.ID,
	}

	cmd.AddGroup(help.GROUP_ORGANIZATIONS)

	cmd.AddCommand(NewListOrganizationsCmd())
	cmd.AddCommand(NewUpdateOrganizationCmd())
	cmd.AddCommand(NewDeleteOrganizationCmd())
	cmd.AddCommand(NewDescribeOrganizationCmd())

	return cmd
}
