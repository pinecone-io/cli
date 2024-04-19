package auth

import (
	"github.com/spf13/cobra"

	logout "github.com/pinecone-io/cli/internal/pkg/auth/logout"
	setApiKey "github.com/pinecone-io/cli/internal/pkg/auth/set_api_key"
)

var helpText = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth <command>",
		Short:   "Authenticate pinecone CLI with your Pinecone account",
		Long: helpText,
	}
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cmd.AddCommand(setApiKey.NewSetApiKeyCmd())
	cmd.AddCommand(logout.NewLogoutCmd())

	return cmd
}

