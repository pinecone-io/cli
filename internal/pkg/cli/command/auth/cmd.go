package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	authHelp = help.Long(`
		Authenticate and manage credentials for the Pinecone CLI.

		There are three ways of authenticating with Pinecone through the CLI:
		through a web browser with user login, using a service account, or 
		configuring a global API key.

		User login and service account authentication provide access to the admin API,
		allowing you to work with projects, API keys, and organizations. You will also
		be able to configure a target organization and project context, which will be
		used for control and data plane operations.
		
		Configuring a global API key overrides any target context, and allows 
		access to control and data plane resources directly. Global API keys
		do not have access to admin API resources.

		See: https://docs.pinecone.io/reference/tools/cli-authentication
	`)
)

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth <command>",
		Short:   "Authenticate and manage credentials for the Pinecone CLI",
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
