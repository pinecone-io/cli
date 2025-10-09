package apiKey

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type describeAPIKeyOptions struct {
	apiKeyID string
	json     bool
}

func NewDescribeAPIKeyCmd() *cobra.Command {
	options := describeAPIKeyOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe an API key by ID",
		Example: help.Examples(`
			pc api-key describe --id "api-key-id"
		`),
		GroupID: help.GROUP_API_KEYS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			apiKey, err := ac.APIKey.Describe(cmd.Context(), options.apiKeyID)
			if err != nil {
				msg.FailMsg("Failed to describe API key %s: %s\n", style.Emphasis(options.apiKeyID), err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(apiKey)
				pcio.Println(json)
			} else {
				presenters.PrintDescribeAPIKeyTable(apiKey)
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.apiKeyID, "id", "i", "", "ID of the API key to describe")
	_ = cmd.MarkFlagRequired("id")

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}
