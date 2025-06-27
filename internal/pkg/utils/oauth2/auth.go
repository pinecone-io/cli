package oauth2

import (
	"context"

	"golang.org/x/oauth2"
)

type Auth struct{}

func (a *Auth) GetAuthURL(ctx context.Context, state string, orgId *string) (string, error) {
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
	opts = append(opts, oauth2.SetAuthURLParam("state", state))
	if orgId != nil && *orgId != "" {
		opts = append(opts, oauth2.SetAuthURLParam("orgId", *orgId))
	}

	return conf.AuthCodeURL("", opts...), nil
}

func (a *Auth) ExchangeAuthCode(ctx context.Context, authCode string) (*oauth2.Token, error) {
	conf, err := newOauth2Config()
	if err != nil {
		return nil, err
	}

	token, err := conf.Exchange(ctx, authCode)
	if err != nil {
		return nil, err
	}
	LogTokenClaims(token, "Obtained access token with auth code")
	return token, nil
}
