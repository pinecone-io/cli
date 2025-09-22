package oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
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
	})
	return mgr, mgrErr
}

func (t *TokenManager) Token(ctx context.Context) (*oauth2.Token, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Pull latest token from secrets
	latestToken := secrets.GetOAuth2Token()

	// If nothing, return the zero token value (callsites check AccessToken)
	if isEmptyToken(&latestToken) {
		var empty oauth2.Token
		return &empty, nil
	}

	// Force refresh slightly before expiry
	if !latestToken.Expiry.IsZero() && time.Until(latestToken.Expiry) <= t.prefetch {
		latestToken.Expiry = time.Now().Add(-time.Minute)
	}

	// Check if we need to refresh
	// - AccessToken is empty but there's a RefreshToken
	// - Inside the prefetch window
	needsRefresh := (latestToken.AccessToken == "" && latestToken.RefreshToken != "") ||
		(!latestToken.Expiry.IsZero() && time.Until(latestToken.Expiry) <= t.prefetch)

	if !needsRefresh {
		return &latestToken, nil
	}

	// Parse the orgId to send along with the refresh request, keeping organization sticky
	claims, err := ParseClaimsUnverified(&latestToken)
	if err != nil {
		return nil, err
	}
	newToken, err := refreshTokenWithOrg(ctx, t.cfg, &latestToken, claims.OrgId)
	if err != nil {
		return nil, err
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

// We need to add orgId to the body in order to keep the organization stick with the refreshed token
// oauth2.TokenSource doesn't do this for us, so this is a custom implementation for calling oauth/token
func refreshTokenWithOrg(ctx context.Context, cfg *oauth2.Config, currToken *oauth2.Token, orgId string) (*oauth2.Token, error) {
	if currToken == nil || currToken.RefreshToken == "" {
		return currToken, nil
	}

	urlForm := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {currToken.RefreshToken},
		"client_id":     {cfg.ClientID},
	}
	if orgId != "" {
		urlForm.Set("orgId", orgId)
	}
	audience, err := getAudience()
	if err != nil {
		return nil, err
	}
	if audience != "" {
		urlForm.Set("audience", audience)
	}

	req, _ := http.NewRequestWithContext(ctx,
		http.MethodPost,
		cfg.Endpoint.TokenURL,
		strings.NewReader(urlForm.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Use the default client unless there's a client passed through Context
	hc := http.DefaultClient
	if v := ctx.Value(oauth2.HTTPClient); v != nil {
		if c, ok := v.(*http.Client); ok && c != nil {
			hc = c
		}
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, pcio.Errorf("failed to refresh user token: %s", resp.Status)
	}

	// Inline struct aligning to what auth0 returns over the wire
	// https://auth0.com/docs/get-started/authentication-and-authorization-flow/authorization-code-flow/call-your-api-using-the-authorization-code-flow
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
		IDToken      string `json:"id_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	if tokenResp.AccessToken == "" {
		return nil, pcio.Errorf("failed to refresh user token: no access token returned")
	}

	// Return a standard oauth2.Token
	// If for some reason TokenType and RefreshToken are empty, use what we've got
	newToken := oauth2.Token{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    firstNonEmpty(tokenResp.TokenType, "Bearer"),
		RefreshToken: firstNonEmpty(tokenResp.RefreshToken, currToken.RefreshToken),
		Expiry:       time.Now().Add(time.Second * time.Duration(tokenResp.ExpiresIn)),
	}
	return &newToken, nil
}

func firstNonEmpty(strs ...string) string {
	for _, str := range strs {
		if str != "" {
			return str
		}
	}
	return ""
}

func isEmptyToken(token *oauth2.Token) bool {
	return token == nil || (token.AccessToken == "" && token.RefreshToken == "")
}

func tokensEqual(t1, t2 *oauth2.Token) bool {
	if t1 == nil || t2 == nil {
		return t1 == t2
	}

	if t1.AccessToken != t2.AccessToken || t1.RefreshToken != t2.RefreshToken {
		return true
	}

	diff := t1.Expiry.Sub(t2.Expiry)
	if diff < 0 {
		diff = -diff
	}
	return diff > 2*time.Second
}
