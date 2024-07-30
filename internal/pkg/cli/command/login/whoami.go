package login

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	pc_oauth2 "github.com/pinecone-io/cli/internal/pkg/utils/oauth2"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewWhoAmICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "whoami",
		Short:   "See the current logged in user",
		GroupID: help.GROUP_START.ID,
		Run: func(cmd *cobra.Command, args []string) {

			accessToken := secrets.OAuth2Token.Get()
			if accessToken.AccessToken == "" {
				msg.InfoMsg("You are not logged in. Please run %s to log in.", style.Code("pinecone login"))
				return
			}

			claims, err := pc_oauth2.ParseClaimsUnverified(&accessToken)
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
