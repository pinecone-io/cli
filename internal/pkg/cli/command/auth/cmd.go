package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	authHelp = help.Long(`
		Authenticate and manage credentials for the Pinecone CLI.
		
		There are three ways to authenticate the CLI with Pinecone:
		1. User login: pc login
		Opens a browser for authentication. Provides full access to Admin API
		(organizations, projects, API keys) and control/data plane operations.

		2. Service account: pc auth configure --client-id "YOUR_CLIENT_ID" --client-secret "YOUR_CLIENT_SECRET"
		Uses client credentials for authentication. Provides the same access
		as user login. Service accounts are created in the Pinecone console.
		
		3. API key: pc auth configure --api-key "YOUR_API_KEY"
		Uses a project API key directly. Provides access to control/data plane
		operations only, but no admin API access.

		See: https://docs.pinecone.io/reference/cli/authentication
	`)
)

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth <command>",
		Short:   "Authenticate and manage credentials for the Pinecone CLI",
		Long:    authHelp,
		GroupID: help.GROUP_AUTH.ID,
	}

	cmd.AddGroup(help.GROUP_AUTH)
	cmd.AddCommand(NewCmdAuthStatus())
	cmd.AddCommand(NewLoginCmd())
	cmd.AddCommand(NewLogoutCmd())
	cmd.AddCommand(NewWhoAmICmd())
	cmd.AddCommand(NewConfigureCmd())
	cmd.AddCommand(NewClearCmd())
	cmd.AddCommand(NewLocalKeysCmd())

	return cmd
}
