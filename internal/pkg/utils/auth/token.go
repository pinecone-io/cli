package auth

import (
	"context"
	"sync"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"golang.org/x/oauth2"
)

const defaultPrefetch = 90 * time.Second

func Token(ctx context.Context) (*oauth2.Token, error) {
	m, err := getTokenManager()
	if err != nil {
		return nil, err
	}
	return m.Token(ctx)
}

type TokenManager struct {
	cfg      *oauth2.Config
	cur      *oauth2.Token
	mu       sync.Mutex
	prefetch time.Duration
}

var (
	mgrOnce sync.Once
	mgr     *TokenManager
	mgrErr  error
)

func getTokenManager() (*TokenManager, error) {
	mgrOnce.Do(func() {
		var cfg *oauth2.Config
		cfg, mgrErr = newOauth2Config()
		if mgrErr != nil {
			return
		}
		mgr = &TokenManager{
			cfg:      cfg,
			prefetch: defaultPrefetch,
		}

		token := secrets.GetOAuth2Token()
		mgr.cur = &token
	})
	return mgr, mgrErr
}

func (t *TokenManager) Token(ctx context.Context) (*oauth2.Token, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Pull latest token from secrets
	latestToken := secrets.GetOAuth2Token()
	if !tokensEqual(t.cur, &latestToken) {
		if isEmptyToken(&latestToken) {
			t.cur = nil
		} else {
			copy := latestToken
			t.cur = &copy
		}
	}

	// If nothing, return the zero token value (callsites check AccessToken)
	if isEmptyToken(t.cur) {
		var empty oauth2.Token
		return &empty, nil
	}

	// Check if we need to refresh
	// - AccessToken is empty but there's a RefreshToken
	// - Inside the prefetch window
	needsRefresh := (t.cur.AccessToken == "" && t.cur.RefreshToken != "") ||
		(!t.cur.Expiry.IsZero() && time.Until(t.cur.Expiry) <= t.prefetch)

	if !needsRefresh {
		return t.cur, nil
	}

	// Force refresh slightly before expiry
	if !t.cur.Expiry.IsZero() && time.Until(t.cur.Expiry) <= t.prefetch {
		t.cur.Expiry = time.Now().Add(-time.Minute)
	}

	base := t.cfg.TokenSource(ctx, t.cur)
	newToken, err := base.Token()
	if err != nil {
		return nil, err
	}

	// Keep refresh token if provider omits on refresh
	if newToken.RefreshToken == "" {
		newToken.RefreshToken = t.cur.RefreshToken
	}

	// Persist token if changed
	if !tokensEqual(t.cur, newToken) {
		secrets.SetOAuth2Token(*newToken)
	}
	t.cur = newToken
	return newToken, nil
}

func (t *TokenManager) ClearToken() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.cur = nil
}

func Logout() {
	secrets.ClearOAuth2Token()

	if mgr != nil {
		mgr.ClearToken()
	}
}

func isEmptyToken(token *oauth2.Token) bool {
	return token == nil || (token.AccessToken == "" && token.RefreshToken == "")
}

func tokensEqual(token1, token2 *oauth2.Token) bool {
	if token1 == nil || token2 == nil {
		return token1 == token2
	}
	return token1.AccessToken == token2.AccessToken &&
		token1.RefreshToken == token2.RefreshToken &&
		token1.Expiry.Equal(token2.Expiry)
}
