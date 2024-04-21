package auth

import (
	"github.com/spf13/cobra"
)

var helpTextAuth = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth <command>",
		Short:   "Authenticate pinecone CLI with your Pinecone account",
		Long: helpTextAuth,
	}
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cmd.AddCommand(NewSetApiKeyCmd())
	cmd.AddCommand(NewLogoutCmd())

	return cmd
}

