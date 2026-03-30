package whoami

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type whoamiCmdOptions struct {
	json bool
}

func NewWhoAmICmd() *cobra.Command {
	options := whoamiCmdOptions{}

	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "See the currently logged in user",
		Example: help.Examples(`
			pc whoami
		`),
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			token, err := oauth.Token(cmd.Context())
			if err != nil {
				msg.FailJSON(options.json, "Error retrieving oauth token: %s", err)
				exit.Error(err, "Error retrieving oauth token")
			}
			if token == nil || token.AccessToken == "" {
				msg.InfoMsg("You are not logged in. Please run %s to log in.", style.Code("pc login"))
				return
			}

			claims, err := oauth.ParseClaimsUnverified(token)
			if err != nil {
				msg.FailJSON(options.json, "An auth token was fetched but an error occurred while parsing the token's claims: %s", err)
				exit.Error(err, "An auth token was fetched but an error occurred while parsing the token's claims")
			}

			if options.json {
				fmt.Println(text.IndentJSON(struct {
					Email string `json:"email"`
					OrgId string `json:"organization_id"`
				}{Email: claims.Email, OrgId: claims.OrgId}))
				return
			}

			msg.InfoMsg("Logged in as %s", style.Emphasis(claims.Email))
		},
	}

	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}
