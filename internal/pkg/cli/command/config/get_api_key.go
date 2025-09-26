package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/spf13/cobra"
)

func NewGetApiKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-api-key",
		Short: "Get the current API key configured for the Pinecone CLI",
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := secrets.GlobalApiKey.Get()
			pcio.Printf("Currently configured global API Key: %s", presenters.MaskHeadTail(apiKey, 4, 4))
		},
	}

	return cmd
}
