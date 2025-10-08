package project

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	projectHelp = help.Long(`
		Manage Pinecone projects. 

		A Pinecone project belongs to an organization and contains a number of 
		resources such as indexes, users, and API keys. Only a user who belongs 
		to the project can access the resources in that project. Each project
		has at least one project owner.
		
		See: https://docs.pinecone.io/guides/projects/understanding-projects
	`)
)

func NewProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project <command>",
		Short:   "Manage Pinecone projects",
		Long:    projectHelp,
		GroupID: help.GROUP_ADMIN.ID,
	}

	cmd.AddGroup(help.GROUP_PROJECTS)

	cmd.AddCommand(NewCreateProjectCmd())
	cmd.AddCommand(NewListProjectsCmd())
	cmd.AddCommand(NewDescribeProjectCmd())
	cmd.AddCommand(NewUpdateProjectCmd())
	cmd.AddCommand(NewDeleteProjectCmd())

	return cmd
}
