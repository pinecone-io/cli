package auth

import (
	_ "embed"
	"io"

	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/login"
	"github.com/spf13/cobra"
)

var (
	loginHelp = help.Long(`
		Log in to the Pinecone CLI using your web browser.

		This is the standard authentication method for interactive use. Logging in
		grants you access to the Admin API (allowing you to manage your
		organizations, projects, and other account-level resources directly
		from the command line), as well as control and data plane operations.

		Running this command opens a browser to the Pinecone login page.
		After you successfully authenticate, the CLI is automatically configured with a
		default target organization and project.
		
		You can view your current target with 'pc target -s' or change it at any
		time with 'pc target -o "ORGANIZATION_NAME" -p "PROECT_NAME"'.
	`)
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Pinecone via user login in a web browser",
		Long:  loginHelp,
		Example: help.Examples(`
			pc auth login
		`),
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
				out = io.Discard
			}

			login.Run(cmd.Context(),
				login.IO{
					In:  cmd.InOrStdin(),
					Out: out,
					Err: cmd.ErrOrStderr(),
				},
				login.Options{},
			)
		},
	}

	return cmd
}
