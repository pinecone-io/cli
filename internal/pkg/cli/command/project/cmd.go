package project

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

func NewProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project <command>",
		Short:   "Manage Pinecone projects",
		GroupID: help.GROUP_MANAGEMENT.ID,
	}

	cmd.AddGroup(help.GROUP_PROJECTS_CRUD)
	cmd.AddGroup(help.GROUP_PROJECTS_API_KEYS)

	cmd.AddCommand(NewListProjectsCmd())
	cmd.AddCommand(NewCreateProjectCmd())
	cmd.AddCommand(NewDeleteProjectCmd())
	cmd.AddCommand(NewListKeysCmd())
	cmd.AddCommand(NewCreateApiKeyCmd())
	cmd.AddCommand(NewDeleteKeyCmd())

	return cmd
}
