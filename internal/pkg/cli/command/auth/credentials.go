package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

func NewCredentialsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "credentials <command>",
		Short:   "Work with managed project credentials",
		GroupID: help.GROUP_AUTH.ID,
	}

	cmd.AddGroup(help.GROUP_AUTH)
	cmd.AddCommand(NewListCredentialsCmd())

	return cmd
}
