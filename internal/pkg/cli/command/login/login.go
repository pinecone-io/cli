package login

import (
	_ "embed"
	"io"

	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/login"
	"github.com/spf13/cobra"
)

var (
	loginHelp = help.Long(`
		Authenticate with Pinecone via user login in a web browser.

		After logging in, a target organization and project context will be automatically set.
		You can set a new target organization or project using pc target before accessing control
		and data plane resources.
	`)
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Pinecone via user login in a web browser",
		Long:  loginHelp,
		Example: help.Examples(`
			pc login
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
