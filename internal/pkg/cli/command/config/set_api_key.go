package config

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewSetApiKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-api-key",
		Short: "Set the API key for the Pinecone CLI",
		Run: func(cmd *cobra.Command, args []string) {
			newApiKey := args[0]
			config.ApiKey.Set(newApiKey)
			fmt.Printf("Config property %s updated. To clear saved keys, run %s.\n", style.Emphasis("api_key"), style.Code("pinecone logout"))
		},
	}

	return cmd
}
