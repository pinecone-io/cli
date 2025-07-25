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

	return &oauth2.Config{
		ClientID: connectionConfig.Auth0ClientId,
		Endpoint: oauth2.Endpoint{
			AuthURL:       connectionConfig.Auth0URL + "/oauth/authorize",
			TokenURL:      connectionConfig.Auth0URL + "/oauth/token",
			DeviceAuthURL: connectionConfig.Auth0URL + "/oauth/device/code",
		},
		Scopes:      []string{"openid", "profile", "email", "offline_access"},
		RedirectURL: "http://127.0.0.1:59049/auth-callback",
	}, nil
}
