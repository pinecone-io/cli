package config

import (
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type GetAPIKeyCmdOptions struct {
	reveal bool
	json   bool
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
			if options.json {
				fmt.Fprintln(os.Stdout, text.IndentJSON(struct {
					APIKey string `json:"api_key"`
				}{APIKey: apiKey}))
				return
			}
			msg.InfoMsg("Current default API key: %s", apiKey)
		},
	}

	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "Reveal the full API key value in the output")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}
