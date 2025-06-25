package login

import (
	"bufio"
	"context"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/browser"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	pc_oauth2 "github.com/pinecone-io/cli/internal/pkg/utils/oauth2"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Login to Pinecone CLI",
		GroupID: help.GROUP_START.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			da := pc_oauth2.DeviceAuth{}
			authResponse, err := da.GetAuthResponse(ctx)
			if err != nil {
				pcio.Println(err)
				return
			}

			pcio.Printf("Visit %s to authorize the CLI.\n", style.Underline(authResponse.VerificationURIComplete))
			pcio.Println()
			pcio.Printf("The code %s should be displayed on the authorization page.\n", style.HeavyEmphasis(authResponse.UserCode))
			pcio.Println()

			// Press enter to launch the browser
			pcio.Printf("Press %s to open the browser.\n", style.Code("[Enter]"))
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			browser.OpenBrowser(authResponse.VerificationURIComplete)

			pcio.Println("After you approve in the browser, it may take a few seconds for the next step to complete.")

			style.Spinner("Waiting for authorization...", func() error {
				token, err := da.GetDeviceAccessToken(ctx, authResponse)
				if err != nil {
					return err
				}

				secrets.OAuth2Token.Set(token)
				return nil
			})

			pcio.Println()
			accessToken := secrets.OAuth2Token.Get()
			claims, err := pc_oauth2.ParseClaimsUnverified(&accessToken)
			if err != nil {
				log.Error().Msg("Error parsing claims")
				msg.FailMsg("An auth token was fetched but an error occurred while parsing the token's claims: %s", err)
				exit.Error(pcio.Errorf("error parsing claims from access token: %s", err))
				return
			}
			msg.SuccessMsg("Logged in as " + style.Emphasis(claims.Email) + ". Defaulted to organization ID: " + style.Emphasis(claims.OrgId))

			// Fetch the user's organizations and projects
			orgsResponse, err := dashboard.ListOrganizations()
			if err != nil {
				log.Error().Msg("Error fetching organizations")
				exit.Error(pcio.Errorf("error fetching organizations: %s", err))
				return
			}

			// target organization is whatever the JWT token's orgId is - defaults on first login currently
			var targetOrg *dashboard.Organization
			for _, org := range orgsResponse.Organizations {
				if org.Id == claims.OrgId {
					targetOrg = &org
					break
				}
			}

			state.TargetOrg.Set(&state.TargetOrganization{
				Name: targetOrg.Name,
				Id:   targetOrg.Id,
			})
			pcio.Println()
			pcio.Printf(style.InfoMsg("Target org set to %s.\n"), style.Emphasis(targetOrg.Name))

			if targetOrg.Projects != nil {
				if len(*targetOrg.Projects) == 0 {
					pcio.Printf(style.InfoMsg("No projects found for organization %s.\n"), style.Emphasis(targetOrg.Name))
					pcio.Println(style.InfoMsg("Please create a project for this organization to work with project resources."))
				} else {
					targetProj := (*targetOrg.Projects)[0]
					state.TargetProj.Set(&state.TargetProject{
						Name: targetProj.Name,
						Id:   targetProj.Id,
					})

					pcio.Printf(style.InfoMsg("Target project set %s.\n"), style.Emphasis(targetProj.Name))
				}
			}

			pcio.Println()
			pcio.Println(style.CodeHint("Run %s to change the target context.", "pinecone target"))

			pcio.Println()
			pcio.Printf("Now try %s to learn about index operations.\n", style.Code("pinecone index -h"))
		},
	}

	return cmd
}
