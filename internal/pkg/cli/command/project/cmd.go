package project

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

func NewProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project <command>",
		Short:   "Manage Pinecone projects",
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
