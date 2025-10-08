package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	localKeysHelp = help.Long(`
		Work with API keys that the CLI is managing in local state.

		When authenticated with user login or a service account, the CLI automatically
		creates and manages API keys for control and data plane operations. This happens
		transparently, the first time you run a control/data plane command ('pc index list').
		You can also create a new key yourself and store it as a managed key using
		'pc api-key create --store'.

		See: https://docs.pinecone.io/reference/tools/cli-authentication
	`)
)

func NewLocalKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "local-keys <command>",
		Short:   "Work with API keys that the CLI is managing in local state",
		Long:    localKeysHelp,
		GroupID: help.GROUP_AUTH.ID,
	}

	cmd.AddGroup(help.GROUP_AUTH)
	cmd.AddCommand(NewListLocalKeysCmd())
	cmd.AddCommand(NewPruneLocalKeysCmd())

	return cmd
}
