package environment

import (
	"fmt"
)

type EnvironmentConnectionSettings struct {
	DashboardUrl   string
	PineconeGCPURL string

	Auth0ClientId string
	Auth0URL      string
	Auth0Audience string
}

var (
	Prod = EnvironmentConnectionSettings{
		DashboardUrl:   "https://console-api.pinecone.io",
		PineconeGCPURL: "https://api.pinecone.io",

		Auth0ClientId: "A4ONXSaOGstwwir0zUztoI6zjyt9zsRH",
		Auth0URL:      "https://login.pinecone.io",
		Auth0Audience: "https://us-central1-production-console.cloudfunctions.net/api/v1",
	}

	Staging = EnvironmentConnectionSettings{
		DashboardUrl:   "https://staging.console-api.pinecone.io",
		PineconeGCPURL: "https://api-staging.pinecone.io",

		Auth0ClientId: "jnuhtpQxTzYw0zrpWdFUEMXS9Bx4FDAR",
		Auth0URL:      "https://internal-beta-pinecone-io.us.auth0.com",
		Auth0Audience: "https://us-central1-console-dev.cloudfunctions.net/api/v1",
	}
)

func GetEnvConfig(env string) (EnvironmentConnectionSettings, error) {
	if env == "production" {
		return Prod, nil
	}

	if env == "staging" {
		return Staging, nil
	}

	return EnvironmentConnectionSettings{}, fmt.Errorf("unknown environment configured: %s", env)
}
