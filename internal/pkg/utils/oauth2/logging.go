package oauth2

import (
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"golang.org/x/oauth2"
)

func LogTokenClaims(token *oauth2.Token, msg string) {
	type MyCustomClaims struct {
		Scope string `json:"scope"`
		Email string `json:"https://pinecone.io/email"`
		jwt.RegisteredClaims
	}

	var p = &jwt.Parser{}
	var claims MyCustomClaims
	p.ParseUnverified(token.AccessToken, &claims)
	log.Debug().
		Str("scope", claims.Scope).
		Str("email", claims.Email).
		Str("sub", claims.Subject).
		Str("iss", claims.Issuer).
		Str("aud", strings.Join(claims.Audience, " ")).
		Str("exp", claims.ExpiresAt.String()).
		Msg(msg)
}
