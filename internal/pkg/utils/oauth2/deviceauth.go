package oauth2

import (
	"context"

	"golang.org/x/oauth2"
)

type DeviceAuth struct{}

func (da *DeviceAuth) GetAuthResponse(ctx context.Context, orgId *string) (*oauth2.DeviceAuthResponse, error) {
	conf, err := newOauth2Config(orgId)
	if err != nil {
		return nil, err
	}

	audience, err := getAudience()
	if err != nil {
		return nil, err
	}

	opts := []oauth2.AuthCodeOption{}
	opts = append(opts, oauth2.SetAuthURLParam("audience", audience))
	if orgId != nil && *orgId != "" {
		opts = append(opts, oauth2.SetAuthURLParam("orgId", *orgId))
	}

	return conf.DeviceAuth(ctx, opts...)
}

func (da *DeviceAuth) GetDeviceAccessToken(ctx context.Context, orgId *string, deviceAuthResponse *oauth2.DeviceAuthResponse) (*oauth2.Token, error) {
	conf, err := newOauth2Config(orgId)
	if err != nil {
		return nil, err
	}
	deviceAuthResponse.Interval += 1 // Add 1 second to the poll interval to avoid slow_down error

	token, err := conf.DeviceAccessToken(ctx, deviceAuthResponse)
	if err != nil {
		return nil, err
	}
	LogTokenClaims(token, "Obtained access token with device auth")
	return token, err
}
