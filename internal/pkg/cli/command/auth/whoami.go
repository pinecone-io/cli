package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewWhoAmICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "See the currently logged in user",
		Example: help.Examples(`
			pc auth whoami
		`),
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {

			token, err := oauth.Token(cmd.Context())
			if err != nil {
				msg.FailMsg("Error retrieving oauth token: %s", err)
				exit.Error().Err(err).Msg("Error retrieving oauth token")
				return
			}
			if token == nil || token.AccessToken == "" {
				msg.InfoMsg("You are not logged in. Please run %s to log in.", style.Code("pc login"))
				return
			}

			claims, err := oauth.ParseClaimsUnverified(token)
			if err != nil {
				msg.FailMsg("An auth token was fetched but an error occurred while parsing the token's claims: %s", err)
				exit.Error().Err(err).Msg("An auth token was fetched but an error occurred while parsing the token's claims")
				return
			}
			msg.InfoMsg("Logged in as " + style.Emphasis(claims.Email))
		},
	}

	return cmd
}
