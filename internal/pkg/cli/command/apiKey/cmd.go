package apiKey

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	apiKeyHelp = help.Long(`
		Work with API keys for a specific Pinecone project. Each Pinecone project has
		one or more API keys. In order to make requests to the Pinecone API, you need 
		to authenticate with an API key.

		In order to work with resources outside of the admin API through the CLI, an 
		API key is required. You can set a global API key using pc auth configure --global-api-key, 
		or store an API key using pc api-key create --store. If you do not explicitly provide a key,
		the CLI will create a managed key for the target project. These API keys can be managed
		using pc auth local-keys.

		See: https://docs.pinecone.io/guides/projects/manage-api-keys
	`)
)

func NewAPIKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "api-key <command>",
		Short:   "Work with API keys for a specific Pinecone project",
		Long:    apiKeyHelp,
		GroupID: help.GROUP_ADMIN.ID,
	}

	cmd.AddGroup(help.GROUP_API_KEYS)
	cmd.AddCommand(NewCreateApiKeyCmd())
	cmd.AddCommand(NewUpdateAPIKeyCmd())
	cmd.AddCommand(NewListKeysCmd())
	cmd.AddCommand(NewDescribeAPIKeyCmd())
	cmd.AddCommand(NewDeleteKeyCmd())

	return cmd
}
