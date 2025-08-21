package apiKey

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

func NewAPIKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "api-key <command>",
		Short:   "Manage API keys for a project",
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
