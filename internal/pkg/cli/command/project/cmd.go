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

	cmd.AddCommand(NewListProjectsCmd())
	cmd.AddCommand(NewCreateProjectCmd())

	return cmd
}
