package oauth2

import (
	"context"
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
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

func GetHttpClient(ctx context.Context, useApiKey bool) *http.Client {
	token := secrets.OAuth2Token.Get()

	if token.AccessToken != "" && !useApiKey {
		log.Debug().Msg("Creating http client with OAuth2 token handling")
		config := newOauth2Config()
		log.Debug().
			Bool("has_access_token", token.AccessToken != "").
			Bool("has_refresh_token", token.AccessToken != "").
			Str("expiry", token.Expiry.String()).
			Msg("Creating http client with OAuth2 token handling")
		LogTokenClaims(&token, "Using saved access token with claims")
		return config.Client(context.Background(), &token)
	}

	log.Debug().Msg("Creating http client without OAuth2 token handling")
	return &http.Client{
		Transport: &apiKeyTransport{
			apiKey: secrets.ApiKey.Get(),
			next:   http.DefaultTransport,
		},
	}
}
