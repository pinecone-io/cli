package auth

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/config"
	"github.com/spf13/cobra"
)

var helpTextSetApiKey = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

func NewSetApiKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-api-key",
		Short: "Set the API key for the Pinecone CLI",
		Long:  helpTextSetApiKey,
		Run: func(cmd *cobra.Command, args []string) {
			newApiKey := args[0]
			config.ApiKey.Set(newApiKey)
			config.SaveConfig()
			fmt.Println("API key set successfully")
		},
	}

	return cmd
}
