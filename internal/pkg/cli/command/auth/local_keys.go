package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

func NewLocalKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "local-keys <command>",
		Short:   "Work with API keys that the CLI is managing locally",
		GroupID: help.GROUP_AUTH.ID,
	}

	cmd.AddGroup(help.GROUP_AUTH)
	cmd.AddCommand(NewListLocalKeysCmd())
	cmd.AddCommand(NewPruneLocalKeysCmd())

	return cmd
}
