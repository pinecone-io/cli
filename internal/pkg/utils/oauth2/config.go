package oauth2

import (
	"golang.org/x/oauth2"
)

func newOauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID: Auth0ClientId,
		Endpoint: oauth2.Endpoint{
			AuthURL:       Auth0URL + "/oauth/authorize",
			TokenURL:      Auth0URL + "/oauth/token",
			DeviceAuthURL: Auth0URL + "/oauth/device/code",
		},
		Scopes:      []string{"openid", "profile", "email", "offline_access"},
		RedirectURL: "http://localhost:59049",
	}
}
