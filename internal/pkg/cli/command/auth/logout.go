package auth

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/config"
	"github.com/spf13/cobra"
)

var helpTextLogout = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Delete all saved credentials from Pinecone CLI configuration",
		Long:  helpTextLogout,
		Run: func(cmd *cobra.Command, args []string) {
			config.ApiKey.Set("")
			config.SaveConfig()
			fmt.Println("Successfully logged out")
		},
	}

	return cmd
}
