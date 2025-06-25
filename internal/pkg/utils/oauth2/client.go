package oauth2

import (
	"context"
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
)

type apiKeyTransport struct {
	apiKey string
	next   http.RoundTripper
}

func (akt *apiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Api-Key", akt.apiKey)
	return akt.next.RoundTrip(req)
}

func GetHttpClient(ctx context.Context, orgId *string) (*http.Client, error) {
	token := secrets.OAuth2Token.Get()
	targetOrgId := state.TargetOrg.Get().Id
	if token.AccessToken != "" {
		log.Debug().Msg("Creating http client with OAuth2 token handling")
		config, err := newOauth2Config(&targetOrgId)
		if err != nil {
			log.Error().Err(err).Msg("Error creating OAuth2 config")
			return nil, err
		}

		log.Debug().
			Bool("has_access_token", token.AccessToken != "").
			Bool("has_refresh_token", token.AccessToken != "").
			Str("expiry", token.Expiry.String()).
			Msg("Creating http client with OAuth2 token handling")
		LogTokenClaims(&token, "Using saved access token with claims")
		return config.Client(context.Background(), &token), nil
	}

	log.Debug().Msg("Creating http client without OAuth2 token handling")
	return &http.Client{
		Transport: &apiKeyTransport{
			apiKey: secrets.ApiKey.Get(),
			next:   http.DefaultTransport,
		},
	}, nil
}
