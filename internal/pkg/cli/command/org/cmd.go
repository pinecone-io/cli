package org

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

func NewOrgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "org <command>",
		Short:   "Manage Pinecone orgs",
		GroupID: help.GROUP_MANAGEMENT.ID,
	}

	cmd.AddCommand(NewListOrgsCmd())

	return cmd
}
