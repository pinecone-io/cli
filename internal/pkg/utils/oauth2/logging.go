package oauth2

import (
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"golang.org/x/oauth2"
)

type MyCustomClaims struct {
	Scope string `json:"scope"`
	Email string `json:"https://pinecone.io/email"`
	OrgId string `json:"https://pinecone.io/orgId"`
	jwt.RegisteredClaims
}

func LogTokenClaims(token *oauth2.Token, msg string) {
	var p = &jwt.Parser{}
	var claims MyCustomClaims
	if token.AccessToken == "" {
		log.Debug().Msg("Cannot parse claims from empty access token")
		return
	}

	p.ParseUnverified(token.AccessToken, &claims)
	exp := "<missing>"
	if claims.ExpiresAt != nil {
		exp = claims.ExpiresAt.String()
	}
	log.Debug().
		Str("scope", claims.Scope).
		Str("email", claims.Email).
		Str("sub", claims.Subject).
		Str("iss", claims.Issuer).
		Str("aud", strings.Join(claims.Audience, " ")).
		Str("exp", exp).
		Msg(msg)
}

func ParseClaimsUnverified(token *oauth2.Token) (*MyCustomClaims, error) {

	var p = &jwt.Parser{}
	var claims MyCustomClaims
	if token.AccessToken == "" {
		return &claims, nil
	}
	_, _, err := p.ParseUnverified(token.AccessToken, &claims)
	if err != nil {
		log.Debug().Err(err).Msg("Error parsing claims from access token")
		exit.Error(pcio.Errorf("error parsing claims from token: %s", err))
	}
	return &claims, err
}
