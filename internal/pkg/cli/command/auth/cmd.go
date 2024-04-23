package auth

import (
	"github.com/spf13/cobra"
)

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate pinecone CLI with your Pinecone account",
	}

	cmd.AddCommand(NewSetApiKeyCmd())
	cmd.AddCommand(NewLogoutCmd())

	return cmd
}
