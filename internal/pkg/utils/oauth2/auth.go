package oauth2

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/oauth2"
)

type Auth struct{}

func (a *Auth) GetAuthURL(ctx context.Context, csrfState string, codeChallenge string, orgId *string) (string, error) {
	conf, err := newOauth2Config()
	if err != nil {
		return "", err
	}

	audience, err := getAudience()
	if err != nil {
		return "", err
	}

	opts := []oauth2.AuthCodeOption{}
	opts = append(opts, oauth2.SetAuthURLParam("audience", audience))
	opts = append(opts, oauth2.SetAuthURLParam("code_challenge", codeChallenge))
	opts = append(opts, oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	if orgId != nil && *orgId != "" {
		opts = append(opts, oauth2.SetAuthURLParam("orgId", *orgId))
	}

	return conf.AuthCodeURL(csrfState, opts...), nil
}

func (a *Auth) ExchangeAuthCode(ctx context.Context, codeVerifier string, authCode string) (*oauth2.Token, error) {
	conf, err := newOauth2Config()
	if err != nil {
		return nil, err
	}

	opts := []oauth2.AuthCodeOption{}
	opts = append(opts, oauth2.SetAuthURLParam("code_verifier", codeVerifier))

	token, err := conf.Exchange(ctx, authCode, opts...)
	if err != nil {
		return nil, err
	}
	LogTokenClaims(token, "Obtained access token with auth code")
	return token, nil
}

func (a *Auth) CreateNewVerifierAndChallenge() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	verifier := base64.RawURLEncoding.EncodeToString(bytes)
	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])
	return verifier, challenge, nil
}
