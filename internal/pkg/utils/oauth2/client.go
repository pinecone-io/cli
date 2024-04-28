package oauth2

import (
	"context"
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
)

func GetHttpClient(ctx context.Context) *http.Client {
	token := secrets.OAuth2Token.Get()
	config := newOauth2Config()
	return config.Client(context.Background(), &token)
}
