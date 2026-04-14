package login

import (
	"bufio"
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"time"

	"golang.org/x/term"

	"github.com/pinecone-io/cli/internal/pkg/utils/browser"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

//go:embed assets/redirect_success.html
var successHTML string

//go:embed assets/redirect_error.html
var errorHTML string

//go:embed assets/pinecone_logo.svg
var logoSVG string

type Options struct {
	Json bool
	// Wait makes the daemon path block until the token is acquired and stored,
	// rather than returning immediately after printing "pending". Use this when
	// the caller needs a valid token on return (e.g. pc target re-auth).
	// RunPostAuthSetup is not called in Wait mode; the caller is responsible
	// for any post-auth state setup and output.
	Wait bool
	// OrgId pins the login flow to a specific organization.
	OrgId *string
	// SSOConnection is the Auth0 connection name to pass as `connection=` in the
	// authorization URL, routing the browser directly to the org's IdP.
	// Callers that hold a valid token before clearing credentials (e.g. pc target)
	// should resolve this with FetchSSOConnection before logout, then pass it here.
	SSOConnection *string
}

func Run(ctx context.Context, opts Options) {
	// Resolve output format once at the top level: explicit --json flag or auto-detected non-TTY stdout.
	opts.Json = opts.Json || !term.IsTerminal(int(os.Stdout.Fd()))

	// In JSON mode, check for a pending session once here. If found, skip the
	// already_authenticated guard and pass the result directly into GetAndSetAccessToken
	// to avoid a redundant directory scan and TOCTOU window.
	if opts.Json {
		sess, result, err := findResumableSession()
		if err != nil {
			msg.FailMsg("Error checking for existing auth session: %s", err)
			exit.Error(err, "Error checking for existing auth session")
		}
		if sess != nil {
			if err := getAndSetAccessTokenJSON(ctx, nil, false, nil, sess, result); err != nil {
				msg.FailMsg("Error acquiring access token while logging in: %s", err)
				exit.Error(err, "Error acquiring access token while logging in")
			}
			return
		}
	}

	token, err := oauth.Token(ctx)

	var te *oauth.TokenError
	expired := errors.As(err, &te) && te.Kind == oauth.TokenErrSessionExpired
	if err != nil && !expired {
		msg.FailMsg("Error retrieving oauth token: %s", err)
		exit.Error(err, "Error retrieving oauth token")
	}

	if !expired && token != nil && token.AccessToken != "" {
		// If --org targets a different organization, re-authenticate now while
		// the token is still valid so we can look up the SSO connection before
		// clearing credentials.
		differentOrg := false
		if opts.OrgId != nil && *opts.OrgId != "" {
			if claims, claimsErr := oauth.ParseClaimsUnverified(token); claimsErr == nil {
				differentOrg = claims.OrgId != *opts.OrgId
			}
		}

		if differentOrg {
			conn, lookupErr := FetchSSOConnection(ctx, *opts.OrgId)
			if lookupErr != nil {
				log.Debug().Err(lookupErr).Msg("SSO connection lookup failed, proceeding without connection param")
			}
			if conn != "" {
				opts.SSOConnection = &conn
			}
			oauth.Logout()
			// Fall through to GetAndSetAccessToken.
		} else {
			// Same org (or no --org flag) — show "already logged in".
			if opts.Json {
				claims, err := oauth.ParseClaimsUnverified(token)
				if err == nil {
					fmt.Fprintln(os.Stdout, text.IndentJSON(struct {
						Status string `json:"status"`
						Email  string `json:"email"`
						OrgId  string `json:"org_id"`
					}{Status: "already_authenticated", Email: claims.Email, OrgId: claims.OrgId}))
				} else {
					fmt.Fprintln(os.Stdout, text.IndentJSON(struct {
						Status string `json:"status"`
					}{Status: "already_authenticated"}))
				}
			} else {
				msg.WarnMsg("You are already logged in. Please log out first using %s.", style.Code("pc auth logout"))
			}
			return
		}
	}

	err = GetAndSetAccessToken(ctx, opts.OrgId, opts)
	if err != nil {
		msg.FailMsg("Error acquiring access token while logging in: %s", err)
		exit.Error(err, "Error acquiring access token while logging in")
	}

	if opts.Json {
		// JSON mode: GetAndSetAccessToken handled all output (pending → polling → authenticated).
		return
	}

	// Non-JSON: complete post-auth setup with human-readable output.
	token, err = oauth.Token(ctx)
	if err != nil {
		msg.FailMsg("Error retrieving oauth token: %s", err)
		exit.Error(err, "Error retrieving oauth token")
	}
	claims, err := oauth.ParseClaimsUnverified(token)
	if err != nil {
		msg.FailMsg("An auth token was fetched but an error occurred while parsing the token's claims: %s", err)
		exit.Error(err, "Error parsing claims from access token")
	}
	msg.Blank()
	msg.SuccessMsg("Logged in as %s. Defaulted to organization ID: %s", style.Emphasis(claims.Email), style.Emphasis(claims.OrgId))

	ac := sdk.NewPineconeAdminClient(ctx)

	orgs, err := ac.Organization.List(ctx)
	if err != nil {
		msg.FailMsg("Error fetching organizations: %s", err)
		exit.Error(err, "Error fetching organizations")
	}

	projects, err := ac.Project.List(ctx)
	if err != nil {
		msg.FailMsg("Error fetching projects: %s", err)
		exit.Error(err, "Error fetching projects")
	}

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

	msg.InfoMsg("Target org set to %s.", style.Emphasis(targetOrg.Name))

	if projects != nil {
		if len(projects) == 0 {
			msg.InfoMsg("No projects found for organization %s.", style.Emphasis(targetOrg.Name))
			msg.InfoMsg("Please create a project for this organization to work with project resources.")
		} else {
			targetProj := projects[0]
			state.TargetProj.Set(state.TargetProject{
				Name: targetProj.Name,
				Id:   targetProj.Id,
			})
			msg.InfoMsg("Target project set %s.", style.Emphasis(targetProj.Name))
		}
	}

	msg.Blank()
	msg.HintMsg("Run %s to change the target context.", style.Code("pc target"))
	msg.HintMsg("Now try %s to learn about index operations.", style.Code("pc index -h"))
}

// GetAndSetAccessToken acquires an OAuth token via the PKCE browser flow and stores it.
//
// Non-JSON mode: runs the callback server inline (blocking), prompts for [Enter] to open
// the browser when stdin is an interactive TTY.
//
// JSON mode: spawns a background daemon that owns the local callback server, prints
// {"status":"pending","url":"..."} immediately so agents can extract the URL from partial
// output, then polls the daemon's result file. When the daemon completes it prints
// {"status":"authenticated",...} and returns. If this process is killed before the daemon
// finishes (e.g. an agent timeout), the daemon keeps running; the next call to this
// function will detect the pending session and resume polling rather than starting a new flow.
func GetAndSetAccessToken(ctx context.Context, orgId *string, opts Options) error {
	// Apply TTY auto-detection here so callers don't have to — if stdout is not
	// a terminal (agentic context), always use the JSON/daemon path.
	opts.Json = opts.Json || !term.IsTerminal(int(os.Stdout.Fd()))
	if opts.Json {
		return getAndSetAccessTokenJSON(ctx, orgId, opts.Wait, opts.SSOConnection, nil, nil)
	}
	return getAndSetAccessTokenInteractive(ctx, orgId, opts.SSOConnection)
}

// getAndSetAccessTokenJSON is the agentic path: daemon-backed, non-blocking on stdin.
// sess and result may be pre-fetched by the caller to avoid a redundant directory scan;
// if nil, findResumableSession is called here.
//
// When wait is false (default for pc login): spawns daemon, prints pending JSON, and
// returns immediately. The caller is expected to call again to poll for the result.
//
// When wait is true (for callers like pc target that need a token on return): spawns
// daemon, blocks until auth completes, and returns with the token stored. RunPostAuthSetup
// is not called; the caller owns post-auth state and output.
func getAndSetAccessTokenJSON(ctx context.Context, orgId *string, wait bool, ssoConnection *string, sess *SessionState, result *SessionResult) error {
	if sess == nil {
		// No pre-fetched session — look one up now.
		var err error
		sess, result, err = findResumableSession()
		if err != nil {
			return fmt.Errorf("error checking for existing auth session: %w", err)
		}
	}
	if sess != nil {
		// If the caller is requesting a specific org that doesn't match the pending
		// session's org, the existing session cannot be used.
		if orgId != nil && sess.OrgId != nil && *orgId != *sess.OrgId {
			if result != nil {
				// Daemon has finished and released the port — clean up and start fresh.
				CleanupSession(sess.SessionId)
				// Fall through to start a new flow.
			} else {
				// Daemon is still running and holds the callback port.
				return fmt.Errorf("an auth session for a different organization (%s) is already in progress; wait for it to expire or complete it first", *sess.OrgId)
			}
		} else {
			return resumeSession(sess, result, wait)
		}
	}

	// Start a new flow.
	a := oauth.Auth{}
	csrfState := randomCSRFState()

	verifier, challenge, err := a.CreateNewVerifierAndChallenge()
	if err != nil {
		return fmt.Errorf("error creating new auth verifier and challenge: %w", err)
	}

	authURL, err := a.GetAuthURL(ctx, csrfState, challenge, orgId, ssoConnection)
	if err != nil {
		return fmt.Errorf("error getting auth URL: %w", err)
	}

	sessionId := newSessionId()
	newSess := &SessionState{
		SessionId: sessionId,
		CSRFState: csrfState,
		AuthURL:   authURL,
		OrgId:     orgId,
		CreatedAt: time.Now(),
	}
	if err := writeSessionState(*newSess); err != nil {
		return fmt.Errorf("error writing session state: %w", err)
	}
	// Pass the PKCE verifier to the daemon via environment variable rather than
	// writing it to the session file, so it never touches disk.
	if err := spawnDaemon(sessionId, verifier); err != nil {
		CleanupSession(sessionId)
		return fmt.Errorf("error spawning auth daemon: %w", err)
	}

	if wait {
		// Caller needs the token on return — block until the daemon completes.
		// Print the auth URL to stderr only: stdout must stay clean so the caller
		// can emit a single JSON document once this function returns.
		fmt.Fprintf(os.Stderr, "Visit the following URL to authenticate:\n\n  %s\n\n", authURL)
		return pollForResult(sessionId, newSess.CreatedAt, true)
	}

	// Agentic login (first call): print pending and return immediately. The daemon
	// owns the callback server and will write a result file when auth completes.
	// The next invocation will detect the pending session and poll for the result.
	printPendingJSON(authURL, sessionId)
	return nil
}

// resumeSession handles a session that was already started (e.g. after a process restart).
// If the daemon already finished, it handles the result immediately. Otherwise it polls.
// When wait is true, RunPostAuthSetup is skipped; the caller owns post-auth state and output.
func resumeSession(sess *SessionState, result *SessionResult, wait bool) error {
	if result != nil {
		// Daemon already completed while we were away.
		defer CleanupSession(sess.SessionId)
		if result.Status == "error" {
			return fmt.Errorf("authentication failed: %s", result.Error)
		}
		_ = secrets.SecretsViper.ReadInConfig()
		if wait {
			return nil
		}
		setupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return RunPostAuthSetup(setupCtx)
	}
	// Still pending — poll until the daemon completes.
	// Don't re-emit pending here: this call will block until done, keeping stdout
	// to a single JSON value per invocation.
	return pollForResult(sess.SessionId, sess.CreatedAt, wait)
}

// pollForResult polls the daemon's result file until auth completes or the session expires.
// On success, when wait is false it calls RunPostAuthSetup to emit authenticated JSON;
// when wait is true it reloads credentials and returns, leaving output to the caller.
// On timeout it returns without cleaning up so the daemon can still complete and be
// resumed on the next invocation.
//
// createdAt is the session's creation time; the deadline is computed from it so that
// the total wait never exceeds sessionMaxAge regardless of when polling started.
//
// The polling loop runs on context.Background() so that the root command's --timeout
// flag (default 60s) does not interrupt a user still authenticating in the browser.
func pollForResult(sessionId string, createdAt time.Time, wait bool) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	remaining := time.Until(createdAt.Add(sessionMaxAge))
	if remaining <= 0 {
		return errors.New("timed out waiting for authentication")
	}
	deadline := time.NewTimer(remaining)
	defer deadline.Stop()

	for {
		select {
		case <-deadline.C:
			// Don't clean up — daemon may still complete; next call will resume.
			return errors.New("timed out waiting for authentication")
		case <-ticker.C:
			result, err := readSessionResult(sessionId)
			if err != nil {
				return fmt.Errorf("error reading session result: %w", err)
			}
			if result == nil {
				continue // daemon not done yet
			}
			defer CleanupSession(sessionId)
			if result.Status == "error" {
				return fmt.Errorf("authentication failed: %s", result.Error)
			}
			// The daemon wrote the token to secrets.yaml from a separate process.
			// Reload from disk so this process's Viper cache reflects it.
			_ = secrets.SecretsViper.ReadInConfig()
			if wait {
				// Caller handles post-auth state and output.
				return nil
			}
			// Use a fresh context for the post-auth API calls: the original ctx may
			// have expired if the user took longer than --timeout to authenticate.
			setupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			return RunPostAuthSetup(setupCtx)
		}
	}
}

func printPendingJSON(authURL, sessionId string) {
	fmt.Fprintln(os.Stdout, text.IndentJSON(struct {
		Status      string `json:"status"`
		URL         string `json:"url"`
		SessionId   string `json:"session_id"`
		Description string `json:"description"`
	}{
		Status:      "pending",
		URL:         authURL,
		SessionId:   sessionId,
		Description: "Navigate to the URL to complete the OAuth authorization flow, then call this command again to retrieve credentials.",
	}))
}

// getAndSetAccessTokenInteractive is the original interactive path: inline callback server,
// optional [Enter]-to-open-browser prompt when stdin is a TTY.
func getAndSetAccessTokenInteractive(ctx context.Context, orgId *string, ssoConnection *string) error {
	// If a daemon from a prior JSON-mode login exists, check whether it has
	// already finished before deciding whether to block interactive login.
	sess, result, err := findResumableSession()
	if err != nil {
		return fmt.Errorf("error checking for existing auth session: %w", err)
	}
	if sess != nil {
		if result != nil {
			// Daemon finished (success or error) and has released the port.
			// Clean up the stale session and proceed with interactive login.
			CleanupSession(sess.SessionId)
		} else {
			// Daemon is still running and holds the callback port.
			fmt.Fprintf(os.Stderr, "An authentication flow is already in progress. Visit the URL below to complete it:\n\n  %s\n\n", sess.AuthURL)
			return fmt.Errorf("authentication already in progress (started at %s) — complete the existing flow or wait for it to expire",
				sess.CreatedAt.Format(time.RFC3339))
		}
	}

	a := oauth.Auth{}
	csrfState := randomCSRFState()

	verifier, challenge, err := a.CreateNewVerifierAndChallenge()
	if err != nil {
		return fmt.Errorf("error creating new auth verifier and challenge: %w", err)
	}

	authURL, err := a.GetAuthURL(ctx, csrfState, challenge, orgId, ssoConnection)
	if err != nil {
		return fmt.Errorf("error getting auth URL: %w", err)
	}

	codeCh := make(chan string, 1)
	serverCtx, cancel := context.WithTimeout(ctx, sessionMaxAge)
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

	fmt.Fprintf(os.Stderr, "Visit %s to authorize the CLI.\n", style.Underline(authURL))

	if term.IsTerminal(int(os.Stdin.Fd())) {
		msg.Blank()
		fmt.Fprintf(os.Stderr, "Press %s to open the browser, or manually paste the URL above.\n", style.Code("[Enter]"))

		go func(ctx context.Context) {
			inputCh := make(chan struct{}, 1)
			go func() {
				_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
				if err != nil {
					log.Error().Err(err).Msg("stdin error: unable to open browser")
					return
				}
				close(inputCh)
			}()
			select {
			case <-ctx.Done():
				return
			case <-inputCh:
				if err := browser.OpenBrowser(authURL); err != nil {
					log.Error().Err(err).Msg("error opening browser")
				}
			case <-time.After(sessionMaxAge):
				return
			}
		}(serverCtx)
	}

	code := <-codeCh
	if code == "" {
		return errors.New("error authenticating CLI and retrieving oauth2 access token")
	}

	token, err := a.ExchangeAuthCode(ctx, verifier, code)
	if err != nil {
		return fmt.Errorf("error exchanging auth code for access token: %w", err)
	}

	claims, err := oauth.ParseClaimsUnverified(token)
	if err != nil {
		log.Error().Err(err).Msg("error parsing claims from access token")
		return err
	}

	if token != nil {
		secrets.SetOAuth2Token(*token)
		secrets.ClientId.Set("")
		secrets.ClientSecret.Set("")

		authContext := state.AuthUserToken
		if state.AuthedUser.Get().AuthContext == state.AuthDefaultAPIKey {
			authContext = state.AuthDefaultAPIKey
		}
		state.AuthedUser.Set(state.TargetUser{
			AuthContext: authContext,
			Email:       claims.Email,
		})
	}

	return nil
}

// applyAuthContext fetches the user's org/project and writes them to state.
// It returns the authenticated user's email so callers can use it without a
// second token fetch. It is the side-effect half of the post-auth setup,
// separated so that callers that don't want to emit JSON (e.g.
// EnsureAuthenticated) can still set context.
func applyAuthContext(ctx context.Context) (email string, err error) {
	token, err := oauth.Token(ctx)
	if err != nil {
		return "", fmt.Errorf("error retrieving oauth token: %w", err)
	}
	claims, err := oauth.ParseClaimsUnverified(token)
	if err != nil {
		return "", fmt.Errorf("error parsing token claims: %w", err)
	}

	ac := sdk.NewPineconeAdminClient(ctx)

	orgs, err := ac.Organization.List(ctx)
	if err != nil {
		return "", fmt.Errorf("error fetching organizations: %w", err)
	}
	projects, err := ac.Project.List(ctx)
	if err != nil {
		return "", fmt.Errorf("error fetching projects: %w", err)
	}

	var targetOrg *pinecone.Organization
	for _, org := range orgs {
		if org.Id == claims.OrgId {
			targetOrg = org
			break
		}
	}
	if targetOrg == nil {
		return "", fmt.Errorf("target organization %s not found", claims.OrgId)
	}

	state.TargetOrg.Set(state.TargetOrganization{
		Name: targetOrg.Name,
		Id:   targetOrg.Id,
	})

	// Always write TargetProj so stale data from a previous session is never
	// returned when the current org has no projects.
	if len(projects) > 0 {
		targetProj := projects[0]
		state.TargetProj.Set(state.TargetProject{
			Name: targetProj.Name,
			Id:   targetProj.Id,
		})
	} else {
		state.TargetProj.Set(state.TargetProject{})
	}

	return claims.Email, nil
}

// RunPostAuthSetup fetches the user's org/project context, sets target defaults,
// and emits the final {"status":"authenticated",...} JSON.
func RunPostAuthSetup(ctx context.Context) error {
	email, err := applyAuthContext(ctx)
	if err != nil {
		return err
	}

	targetOrg := state.TargetOrg.Get()
	targetProj := state.TargetProj.Get()

	fmt.Fprintln(os.Stdout, text.IndentJSON(struct {
		Status      string `json:"status"`
		Email       string `json:"email"`
		OrgId       string `json:"org_id"`
		OrgName     string `json:"org_name"`
		ProjectId   string `json:"project_id"`
		ProjectName string `json:"project_name"`
	}{Status: "authenticated", Email: email, OrgId: targetOrg.Id, OrgName: targetOrg.Name, ProjectId: targetProj.Id, ProjectName: targetProj.Name}))

	return nil
}

// spawnDaemon starts a detached `pc auth _daemon --session-id <id>` process.
// The PKCE verifier is passed via environment variable so it never touches disk.
func spawnDaemon(sessionId, pkceVerifier string) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error finding executable path: %w", err)
	}
	cmd := exec.Command(exe, "auth", "_daemon", "--session-id", sessionId)
	detachProcess(cmd)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Env = append(os.Environ(), "PINECONE_PKCE_VERIFIER="+pkceVerifier)
	return cmd.Start()
}

func ServeAuthCodeListener(ctx context.Context, csrfState string) (string, error) {
	codeCh := make(chan string)
	errCh := make(chan error)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth-callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if state != csrfState {
			errCh <- fmt.Errorf("state mismatch on authentication")
			return
		}

		templateData := map[string]template.HTML{"LogoSVG": template.HTML(logoSVG)}
		if code == "" {
			if err := renderHTML(w, errorHTML, templateData); err != nil {
				errCh <- fmt.Errorf("error rendering authentication response HTML: %w", err)
				return
			}
		} else {
			if err := renderHTML(w, successHTML, templateData); err != nil {
				errCh <- fmt.Errorf("error rendering authentication response HTML: %w", err)
				return
			}
		}
		w.(http.Flusher).Flush()
		codeCh <- code
	})

	serve := &http.Server{
		Addr:    "127.0.0.1:59049",
		Handler: mux,
	}
	go func() {
		if err := serve.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("error listening for auth code: %w", err)
			return
		}
	}()

	select {
	case code := <-codeCh:
		_ = serve.Shutdown(ctx)
		return code, nil
	case err := <-errCh:
		_ = serve.Shutdown(ctx)
		return "", err
	case <-ctx.Done():
		_ = serve.Shutdown(ctx)
		if ctx.Err() != nil {
			return "", fmt.Errorf("error waiting for authorization: %w", ctx.Err())
		}
	}

	return "", errors.New("error waiting for authentication response")
}

func renderHTML(w http.ResponseWriter, htmlTemplate string, data map[string]template.HTML) error {
	tmpl, err := template.New("auth-response").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("error parsing auth response HTML template: %w", err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("error executing auth response HTML template: %w", err)
	}
	return nil
}

// EnsureAuthenticated verifies that valid credentials are available, transparently
// completing a finished pending session if one exists.
//
//   - Valid token or non-OAuth credentials present → returns nil immediately.
//   - Pending session whose daemon has already completed → reloads credentials from
//     disk and returns nil. This is the "lazy completion" path: after a successful
//     browser login the next command just works without a second `pc login` call.
//   - Pending session still in progress → returns an error containing the auth URL.
//   - No credentials and no session → returns a "not authenticated" error.
func EnsureAuthenticated(ctx context.Context) error {
	// Service-account and API key credentials don't use OAuth tokens.
	if secrets.ClientId.Get() != "" && secrets.ClientSecret.Get() != "" {
		return nil
	}
	if secrets.DefaultAPIKey.Get() != "" {
		return nil
	}

	// Check for a valid OAuth token.
	token, err := oauth.Token(ctx)
	var te *oauth.TokenError
	if err != nil {
		if !errors.As(err, &te) || te.Kind != oauth.TokenErrSessionExpired {
			// Real error (network failure, etc.) — propagate it.
			return err
		}
		// Session expired — fall through to check for a pending session.
	} else if token != nil && token.AccessToken != "" {
		// Valid token. If no target org is set yet (e.g. first command after a fresh
		// login before a second `pc login` call was made), initialize the context now
		// so the calling command doesn't fail with "need to target a project".
		if state.TargetOrg.Get().Id == "" {
			if _, err := applyAuthContext(ctx); err != nil {
				log.Debug().Err(err).Msg("EnsureAuthenticated: applyAuthContext failed")
			}
		}
		return nil
	}

	// No valid token — look for a pending session whose daemon may have finished.
	sess, result, sessErr := findResumableSession()
	if sessErr != nil {
		return fmt.Errorf("error checking auth session: %w", sessErr)
	}

	if sess == nil {
		if err != nil {
			// Token was expired and there's no pending session to fall back to.
			return err
		}
		return fmt.Errorf("not authenticated. Run %s to log in.", style.Code("pc login"))
	}

	if result == nil {
		// Daemon still running — auth not yet complete.
		return fmt.Errorf("authentication in progress. Visit the following URL to complete login, then retry:\n\n  %s\n\nOr run %s to check status.", sess.AuthURL, style.Code("pc login -j"))
	}

	// Daemon finished.
	defer CleanupSession(sess.SessionId)
	if result.Status == "error" {
		return fmt.Errorf("authentication failed: %s. Run %s to try again.", result.Error, style.Code("pc login"))
	}

	// Reload credentials written by the daemon process into this process's cache,
	// then set the target org/project context so the calling command can proceed
	// without a separate `pc login` or `pc target` call.
	_ = secrets.SecretsViper.ReadInConfig()
	if _, err := applyAuthContext(ctx); err != nil {
		// Non-fatal: credentials are valid, context setup is best-effort.
		log.Debug().Err(err).Msg("EnsureAuthenticated: applyAuthContext failed after lazy credential reload")
	}
	return nil
}

func randomCSRFState() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
