package project

import (
	"github.com/spf13/cobra"
)

func NewProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project <command>",
		Short: "Manage Pinecone projects",
	}

	cmd.AddCommand(NewListProjectsCmd())

	return cmd
}
