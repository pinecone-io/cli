package deviceauth

import (
	"context"

	"golang.org/x/oauth2"
)

// These are set with ldflags during build process
var Auth0ClientId = "XY4m3uRYoW6S0dK9ypXmM0Wc2bzUAdXW"
var Auth0URL = "https://login.pinecone.io"
var Auth0Audience = "https://us-central1-production-console.cloudfunctions.net/api/v1"

type DeviceAuth struct{}

func newOauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID: Auth0ClientId,
		Endpoint: oauth2.Endpoint{
			AuthURL:       Auth0URL + "/oauth/authorize",
			TokenURL:      Auth0URL + "/oauth/token",
			DeviceAuthURL: Auth0URL + "/oauth/device/code",
		},
		Scopes:      []string{"openid", "profile", "email"},
		RedirectURL: "http://localhost:59049",
	}
}

func (da *DeviceAuth) GetAuthResponse(ctx context.Context) (*oauth2.DeviceAuthResponse, error) {
	conf := newOauth2Config()
	opts := oauth2.SetAuthURLParam("audience", Auth0Audience)

	return conf.DeviceAuth(ctx, opts)
}

func (da *DeviceAuth) GetDeviceAccessToken(ctx context.Context, deviceAuthResponse *oauth2.DeviceAuthResponse) (*oauth2.Token, error) {
	conf := newOauth2Config()
	deviceAuthResponse.Interval += 1 // Add 1 second to the poll interval to avoid slow_down error
	token, err := conf.DeviceAccessToken(ctx, deviceAuthResponse)
	return token, err
}
