package oauth2

import (
	"context"
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
)

func GetHttpClient(ctx context.Context) *http.Client {
	token := secrets.OAuth2Token.Get()
	config := newOauth2Config()
	log.Debug().
		Bool("has_access_token", token.AccessToken != "").
		Bool("has_refresh_token", token.AccessToken != "").
		Str("expiry", token.Expiry.String()).
		Msg("Creating http client with OAuth2 token handling")
	LogTokenClaims(&token, "Using saved access token with claims")
	return config.Client(context.Background(), &token)
}
