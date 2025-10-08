package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/spf13/cobra"
)

type GetAPIKeyCmdOptions struct {
	reveal bool
}

func NewGetApiKeyCmd() *cobra.Command {
	options := GetAPIKeyCmdOptions{}

	cmd := &cobra.Command{
		Use:   "get-api-key",
		Short: "Get the current default API key configured for the Pinecone CLI",
		Example: help.Examples(`
		    pc config get-api-key
		`),
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := secrets.DefaultAPIKey.Get()
			if !options.reveal {
				apiKey = presenters.MaskHeadTail(apiKey, 4, 4)
			}
			pcio.Printf("Current default API key: %s", apiKey)
		},
	}

	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "Reveal the full API key value in the output")

	return cmd
}
