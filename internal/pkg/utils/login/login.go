package login

import (
	"bufio"
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/browser"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

//go:embed assets/redirect_success.html
var successHTML string

//go:embed assets/redirect_error.html
var errorHTML string

//go:embed assets/pinecone_logo.svg
var logoSVG string

type IO struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type Options struct{}

func Run(ctx context.Context, io IO, opts Options) {
	err := GetAndSetAccessToken(nil)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching authentication token")
		exit.Error(pcio.Errorf("error acquiring access token while logging in: %w", err))
	}

	// Parse token claims to get orgId
	accessToken, err := oauth.Token(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving oauth token")
		msg.FailMsg("Error retrieving oauth token: %s", err)
		exit.Error(pcio.Errorf("error retrieving oauth token: %w", err))
	}
	claims, err := oauth.ParseClaimsUnverified(accessToken)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing authentication token claims")
		msg.FailMsg("An auth token was fetched but an error occurred while parsing the token's claims: %s", err)
		exit.Error(pcio.Errorf("error parsing claims from access token: %w", err))
	}
	msg.SuccessMsg("Logged in as " + style.Emphasis(claims.Email) + ". Defaulted to organization ID: " + style.Emphasis(claims.OrgId))

	ac := sdk.NewPineconeAdminClient()
	if err != nil {
		log.Error().Err(err).Msg("Error creating Pinecone admin client")
		exit.Error(pcio.Errorf("error creating Pinecone admin client: %w", err))
	}

	// Fetch the user's organizations and projects for the default org associated with the JWT token (if it exists)
	orgs, err := ac.Organization.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching organizations")
		exit.Error(pcio.Errorf("error fetching organizations: %w", err))
	}

	projects, err := ac.Project.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching projects")
		exit.Error(pcio.Errorf("error fetching projects: %w", err))
	}

	// target organization is whatever the JWT token's orgId is - defaults on first login currently
	var targetOrg *pinecone.Organization
	for _, org := range orgs {
		if org.Id == claims.OrgId {
			targetOrg = org
			break
		}
	}

	state.TargetOrg.Set(state.TargetOrganization{
		Name: targetOrg.Name,
		Id:   targetOrg.Id,
	})
	pcio.Println()
	pcio.Printf(style.InfoMsg("Target org set to %s.\n"), style.Emphasis(targetOrg.Name))

	if projects != nil {
		if len(projects) == 0 {
			pcio.Printf(style.InfoMsg("No projects found for organization %s.\n"), style.Emphasis(targetOrg.Name))
			pcio.Println(style.InfoMsg("Please create a project for this organization to work with project resources."))
		} else {
			targetProj := projects[0]
			state.TargetProj.Set(state.TargetProject{
				Name: targetProj.Name,
				Id:   targetProj.Id,
			})

			pcio.Printf(style.InfoMsg("Target project set %s.\n"), style.Emphasis(targetProj.Name))
		}
	}

	pcio.Println()
	pcio.Println(style.CodeHint("Run %s to change the target context.", style.Code("pc target")))

	pcio.Println()
	pcio.Printf("Now try %s to learn about index operations.\n", style.Code("pc index -h"))
}

// Takes an optional orgId, and attempts to acquire an access token scoped to the orgId if provided.
// If a token is successfully acquired it's set in the secrets store, and the user is considered logged in with state.AuthUserToken.
func GetAndSetAccessToken(orgId *string) error {
	ctx := context.Background()
	a := oauth.Auth{}

	// CSRF state
	csrfState := randomCSRFState()

	// PKCE verifier and challenge
	verifier, challenge, err := a.CreateNewVerifierAndChallenge()
	if err != nil {
		exit.Error(pcio.Error("error creating new auth verifier and challenge"))
		return err
	}

	authURL, err := a.GetAuthURL(ctx, csrfState, challenge, orgId)
	if err != nil {
		exit.Error(pcio.Errorf("error getting auth URL: %w", err))
		return err
	}

	// Spin up a local server in a goroutine to handle receiving the authorization code from auth0
	codeCh := make(chan string, 1)
	serverCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	go func() {
		code, err := ServeAuthCodeListener(serverCtx, csrfState)
		if err != nil {
			log.Error().Err(err).Msg("Error waiting for authorization")
			codeCh <- ""
			return
		}
		codeCh <- code
	}()

	pcio.Printf("Visit %s to authorize the CLI.\n", style.Underline(authURL))
	pcio.Println()
	pcio.Printf("Press %s to open the browser, or manually paste the URL above.\n", style.Code("[Enter]"))

	// spawn a goroutine to optionally wait for [Enter] as input
	go func(ctx context.Context) {
		// inner channel to signal that [Enter] was pressed
		inputCh := make(chan struct{}, 1)

		// spawn inner goroutine to read stdin (blocking)
		go func() {
			_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
			if err != nil {
				log.Error().Err(err).Msg("stdin error: unable to open browser")
				return
			}
			close(inputCh)
		}()

		// wait for [Enter], auth code, or timeout
		select {
		case <-ctx.Done():
			return
		case <-inputCh:
			err = browser.OpenBrowser(authURL)
			if err != nil {
				log.Error().Err(err).Msg("error opening browser")
				return
			}
		case <-time.After(5 * time.Minute):
			// extra precaution to prevent hanging indefinitely on stdin
			return
		}
	}(serverCtx)

	// Wait for auth code and exchange for access token
	code := <-codeCh
	if code == "" {
		return pcio.Error("error authenticating CLI and retrieving oauth2 access token")
	}

	token, err := a.ExchangeAuthCode(ctx, verifier, code)
	if err != nil {
		exit.Error(pcio.Errorf("error exchanging auth code for access token: %w", err))
	}

	claims, err := oauth.ParseClaimsUnverified(token)
	if err != nil {
		log.Error().Err(err).Msg("error parsing claims from access token")
		return err
	}

	if token != nil {
		// Store the token
		secrets.SetOAuth2Token(*token)

		// Clear any existing client_id and client_secret values
		secrets.ClientId.Set("")
		secrets.ClientSecret.Set("")

		// Update target credentials context
		// TODO - proper getters / setters for state
		globalAPIKey := secrets.GlobalApiKey.Get()
		authContext := state.AuthUserToken
		if state.TargetCreds.Get().AuthContext == state.AuthGlobalAPIKey {
			authContext = state.AuthGlobalAPIKey
		}
		state.TargetCreds.Set(state.TargetUser{
			AuthContext:  authContext,
			Email:        claims.Email,
			GlobalAPIKey: globalAPIKey,
		})
	}

	return nil
}

func ServeAuthCodeListener(ctx context.Context, csrfState string) (string, error) {
	codeCh := make(chan string)

	// start server to receive the auth code
	mux := http.NewServeMux()
	mux.HandleFunc("/auth-callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if state != csrfState {
			exit.Error(pcio.Errorf("state mismatch on authentication"))
			return
		}

		// Code is empty, there was an error authenticating, return error HTML
		templateData := map[string]template.HTML{"LogoSVG": template.HTML(logoSVG)}
		if code == "" {
			if err := renderHTML(w, errorHTML, templateData); err != nil {
				exit.Error(pcio.Errorf("error rendering authentication response HTML: %w", err))
				return
			}
		} else {
			if err := renderHTML(w, successHTML, templateData); err != nil {
				exit.Error(pcio.Errorf("error rendering authentication response HTML: %w", err))
				return
			}
		}
		w.(http.Flusher).Flush()
		codeCh <- code
	})

	// Start server and listen for auth code response
	serve := &http.Server{
		Addr:    "127.0.0.1:59049",
		Handler: mux,
	}
	go func() {
		if err := serve.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			exit.Error(pcio.Errorf("error listening for auth code: %w", err))
			return
		}
	}()

	select {
	case code := <-codeCh:
		_ = serve.Shutdown(ctx)
		return code, nil
	case <-ctx.Done():
		_ = serve.Shutdown(ctx)
		if ctx.Err() != nil {
			exit.Error(pcio.Errorf("error waiting for authorization: %w", ctx.Err()))
			return "", ctx.Err()
		}
	}

	return "", pcio.Error("error waiting for authentication response")
}

func renderHTML(w http.ResponseWriter, htmlTemplate string, data map[string]template.HTML) error {
	tmpl, err := template.New("auth-response").Parse(htmlTemplate)
	if err != nil {
		exit.Error(pcio.Errorf("error parsing auth response HTML template: %w", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, data); err != nil {
		exit.Error(pcio.Errorf("error executing auth response HTML template: %w", err))
		return err
	}
	return nil
}

func randomCSRFState() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
