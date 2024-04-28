package oauth2

import (
	"context"
	"net/http"

	gooauth2 "golang.org/x/oauth2"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
)

func GetHttpClient(ctx context.Context) *http.Client {
	token := gooauth2.Token{
		AccessToken: secrets.AccessToken.Get(),
	}
	config := newOauth2Config()
	return config.Client(context.Background(), &token)
}
