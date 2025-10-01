package auth

import (
	_ "embed"
	"io"

	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/login"
	"github.com/spf13/cobra"
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to the Pinecone CLI through the browser using your Pinecone account",
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
