package oauth2

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/environment"
	"golang.org/x/oauth2"
)

func getAudience() (string, error) {
	connectionConfig, err := environment.GetEnvConfig(config.Environment.Get())
	if err != nil {
		return "", err
	}

	return connectionConfig.Auth0Audience, nil
}

func newOauth2Config() (*oauth2.Config, error) {
	connectionConfig, err := environment.GetEnvConfig(config.Environment.Get())
	if err != nil {
		return nil, err
	}

	// TODO: figure out if we need to actually modify these urls
	// authURLPath := "/oauth/authorize"
	// deviceAuthURLPath := "/oauth/device/code"
	// if orgId != nil && *orgId != "" {
	// 	authURLPath = authURLPath + "orgId=" + *orgId
	// 	deviceAuthURLPath = deviceAuthURLPath + "orgId=" + *orgId
	// }

	return &oauth2.Config{
		ClientID: connectionConfig.Auth0ClientId,
		Endpoint: oauth2.Endpoint{
			AuthURL:       connectionConfig.Auth0URL + "/oauth/authorize",
			TokenURL:      connectionConfig.Auth0URL + "/oauth/token",
			DeviceAuthURL: connectionConfig.Auth0URL + "/oauth/device/code",
		},
		Scopes:      []string{"openid", "profile", "email", "offline_access"},
		RedirectURL: "http://localhost:59049",
	}, nil
}
