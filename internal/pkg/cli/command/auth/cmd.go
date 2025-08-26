package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth <command>",
		Short:   "Log in and manage authentication credentials for the CLI",
		GroupID: help.GROUP_AUTH.ID,
	}

	cmd.AddGroup(help.GROUP_AUTH)
	cmd.AddCommand(NewLoginCmd())
	cmd.AddCommand(NewLogoutCmd())
	cmd.AddCommand(NewWhoAmICmd())
	cmd.AddCommand(NewConfigureCmd())

	return cmd
}
