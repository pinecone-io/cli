package login

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
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
			err := GetAndSetAccessToken(nil)
			if err != nil {
				exit.Error(pcio.Errorf("error acquiring access token while logging in: %w", err))
			}

			// Parse token claims to get orgId
			accessToken := secrets.OAuth2Token.Get()
			claims, err := pc_oauth2.ParseClaimsUnverified(&accessToken)
			if err != nil {
				log.Error().Msg("Error parsing claims")
				msg.FailMsg("An auth token was fetched but an error occurred while parsing the token's claims: %s", err)
				exit.Error(pcio.Errorf("error parsing claims from access token: %w", err))
				return
			}
			msg.SuccessMsg("Logged in as " + style.Emphasis(claims.Email) + ". Defaulted to organization ID: " + style.Emphasis(claims.OrgId))

			// Fetch the user's organizations and projects for the default org associated with the JWT token (if it exists)
			orgsResponse, err := dashboard.ListOrganizations()
			if err != nil {
				log.Error().Msg("Error fetching organizations")
				exit.Error(pcio.Errorf("error fetching organizations: %w", err))
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

// Takes an optional orgId, and attempts to acquire an access token scoped to the orgId if provided
// If a token is successfully acquired it's set in the secrets store, and the user is considered logged in
func GetAndSetAccessToken(orgId *string) error {
	ctx := context.Background()

	// da := pc_oauth2.DeviceAuth{}
	a := pc_oauth2.Auth{}

	// CSRF state
	csrfState := randomState()

	// PKCE verifier and challenge
	verifier, challenge, err := a.CreateNewVerifierAndChallenge()

	authURL, err := a.GetAuthURL(ctx, csrfState, challenge, orgId)
	if err != nil {
		exit.Error(pcio.Errorf("error getting auth URL: %w", err))
		return err
	}

	pcio.Printf("Visit %s to authorize the CLI.\n", style.Underline(authURL))
	pcio.Println()

	pcio.Printf("Press %s to open the browser.\n", style.Code("[Enter]"))
	_, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		exit.Error(pcio.Errorf("error reading input: %w", err))
		return err
	}

	err = browser.OpenBrowser(authURL)
	if err != nil {
		exit.Error(pcio.Errorf("error opening browser: %w", err))
		return err
	}
	pcio.Println("After you approve in the browser, it may take a few seconds for the next step to complete.")

	// Spin up a local server to handle receiving the authorization code from auth0
	var code string
	err = style.Spinner("Waiting for authorization...\n", func() error {
		codeCh := make(chan string)

		// start server to receive the auth code
		mux := http.NewServeMux()
		mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			state := r.URL.Query().Get("state")

			if state != csrfState {
				exit.Error(pcio.Errorf("State mismatch on authentication"))
				return
			}
			codeCh <- code
		})
		serve := &http.Server{Addr: "127.0.0.1:59049", Handler: mux}

		go func() {
			if err := serve.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				exit.Error(pcio.Errorf("error listening for auth code: %w", err))
				return
			}
		}()

		select {
		case code = <-codeCh:
		case <-ctx.Done():
			_ = serve.Shutdown(ctx)
			if ctx.Err() != nil {
				exit.Error(pcio.Errorf("error waiting for authorization: %w", ctx.Err()))
				return ctx.Err()
			}
		}

		_ = serve.Shutdown(ctx)
		return nil
	})
	if err != nil {
		exit.Error(pcio.Errorf("error waiting for authorization: %w", err))
		return err
	}

	// Exchange auth code for access token
	if code != "" {
		token, err := a.ExchangeAuthCode(ctx, verifier, code)
		if err != nil {
			exit.Error(pcio.Errorf("error exchanging auth code for access token: %w", err))
			return err
		}
		secrets.OAuth2Token.Set(token)
		return nil
	}

	// if we're here, it means we were not able to authenticate with the auth0 server
	exit.Error(pcio.Errorf("error authenticating CLI and retrieving oauth2 access token"))
	return pcio.Errorf("error authenticating CLI and retrieving oauth2 access token")
}

func randomState() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
