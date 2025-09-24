package auth

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewWhoAmICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "See the current logged in user",
		Example: heredoc.Doc(`
		$ pc auth whoami
		`),
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {

			accessToken, err := oauth.Token(cmd.Context())
			if err != nil {
				log.Error().Err(err).Msg("Error retrieving oauth token")
				msg.FailMsg("Error retrieving oauth token: %s", err)
				exit.Error(pcio.Errorf("error retrieving oauth token: %w", err))
				return
			}
			if accessToken.AccessToken == "" {
				msg.InfoMsg("You are not logged in. Please run %s to log in.", style.Code("pc login"))
				return
			}

			claims, err := oauth.ParseClaimsUnverified(accessToken)
			if err != nil {
				log.Error().Msg("Error parsing claims")
				msg.FailMsg("An auth token was fetched but an error occurred while parsing the token's claims: %s", err)
				exit.Error(pcio.Errorf("error parsing claims from access token: %s", err))
				return
			}
			msg.InfoMsg("Logged in as " + style.Emphasis(claims.Email))
		},
	}

	return cmd
}
