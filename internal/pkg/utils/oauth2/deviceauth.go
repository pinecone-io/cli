package oauth2

import (
	"context"

	"golang.org/x/oauth2"
)

type DeviceAuth struct{}

func (da *DeviceAuth) GetAuthResponse(ctx context.Context) (*oauth2.DeviceAuthResponse, error) {
	conf := newOauth2Config()
	opts := oauth2.SetAuthURLParam("audience", Auth0Audience)

	return conf.DeviceAuth(ctx, opts)
}

func (da *DeviceAuth) GetDeviceAccessToken(ctx context.Context, deviceAuthResponse *oauth2.DeviceAuthResponse) (*oauth2.Token, error) {
	conf := newOauth2Config()
	deviceAuthResponse.Interval += 1 // Add 1 second to the poll interval to avoid slow_down error

	token, err := conf.DeviceAccessToken(ctx, deviceAuthResponse, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, err
	}
	LogTokenClaims(token, "Obtained access token with device auth")
	return token, err
}
