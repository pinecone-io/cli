package login

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
)

// RunDaemon is called by the hidden `pc auth _daemon` subcommand. It reads the
// session state written by the parent, runs the OAuth callback server, exchanges
// the auth code for a token, stores the token, and writes a result file.
// It must not write to stdout/stderr since it runs detached from any terminal.
func RunDaemon(sessionId string) {
	ctx, cancel := context.WithTimeout(context.Background(), sessionMaxAge)
	defer cancel()

	sess, err := ReadSessionState(sessionId)
	if err != nil {
		writeDaemonError(sessionId, fmt.Sprintf("error reading session state: %s", err))
		return
	}

	codeCh := make(chan string, 1)
	go func() {
		code, err := ServeAuthCodeListener(ctx, sess.CSRFState)
		if err != nil {
			log.Error().Err(err).Str("session_id", sessionId).Msg("daemon: error waiting for authorization")
			codeCh <- ""
			return
		}
		codeCh <- code
	}()

	code := <-codeCh
	if code == "" {
		writeDaemonError(sessionId, "authentication timed out or failed: no auth code received")
		return
	}

	pkceVerifier := os.Getenv("PINECONE_PKCE_VERIFIER")
	if pkceVerifier == "" {
		writeDaemonError(sessionId, "PKCE verifier not provided to daemon")
		return
	}

	a := oauth.Auth{}
	token, err := a.ExchangeAuthCode(ctx, pkceVerifier, code)
	if err != nil {
		writeDaemonError(sessionId, fmt.Sprintf("error exchanging auth code: %s", err))
		return
	}
	if token == nil {
		writeDaemonError(sessionId, "error exchanging auth code: no token returned")
		return
	}

	claims, err := oauth.ParseClaimsUnverified(token)
	if err != nil {
		writeDaemonError(sessionId, fmt.Sprintf("error parsing token claims: %s", err))
		return
	}

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

	_ = WriteSessionResult(SessionResult{
		SessionId:   sessionId,
		Status:      "success",
		CompletedAt: time.Now(),
	})
}

func writeDaemonError(sessionId, errMsg string) {
	_ = WriteSessionResult(SessionResult{
		SessionId:   sessionId,
		Status:      "error",
		Error:       errMsg,
		CompletedAt: time.Now(),
	})
}
