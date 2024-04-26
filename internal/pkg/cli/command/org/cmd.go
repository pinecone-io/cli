package org

import (
	"github.com/spf13/cobra"
)

func NewOrgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org <command>",
		Short: "Manage Pinecone orgs",
	}

	cmd.AddCommand(NewListOrgsCmd())

	return cmd
}
