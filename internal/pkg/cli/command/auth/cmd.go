package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/cli/command/config"
	"github.com/spf13/cobra"
)

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate pinecone CLI with your Pinecone account",
	}

	cmd.AddCommand(config.NewSetApiKeyCmd())
	cmd.AddCommand(NewLogoutCmd())
	cmd.AddCommand(NewLoginCmd())

	return cmd
}
