package apiKey

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	apiKeyHelp = help.Long(`
		Work with API keys for a Pinecone project.

		API keys are used to authenticate with the Pinecone API. You can set a default 
		API key using 'pc auth configure --api-key', or you can create and store a 
		new one for the current project with 'pc api-key create --store'. 
		
		If you do not provide a key or store one, the CLI creates a "managed key" for the project.
		These keys can be viewed and managed with 'pc auth local-keys'.

		See: https://docs.pinecone.io/guides/projects/manage-api-keys
	`)
)

func NewAPIKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "api-key <command>",
		Short:   "Work with API keys for a Pinecone project",
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
