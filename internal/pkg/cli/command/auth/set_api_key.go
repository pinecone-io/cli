package auth

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/config"
	"github.com/spf13/cobra"
)

func NewSetApiKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-api-key",
		Short: "Set the API key for the Pinecone CLI",
		Run: func(cmd *cobra.Command, args []string) {
			newApiKey := args[0]
			config.ApiKey.Set(newApiKey)
			config.SaveConfig()
			fmt.Println("API key set successfully")
		},
	}

	return cmd
}
