package environment

import (
	"fmt"
)

type EnvironmentConnectionSettings struct {
	DashboardUrl             string
	IndexControlPlaneUrl     string
	AssistantControlPlaneUrl string
	AssistantDataPlaneUrl    string

	Auth0ClientId string
	Auth0URL      string
	Auth0Audience string
}

var (
	Prod = EnvironmentConnectionSettings{
		DashboardUrl:             "https://console-api.pinecone.io",
		IndexControlPlaneUrl:     "https://api.pinecone.io",
		AssistantControlPlaneUrl: "https://api.pinecone.io",
		AssistantDataPlaneUrl:    "https://prod-1-data.ke.pinecone.io",

		Auth0ClientId: "A4ONXSaOGstwwir0zUztoI6zjyt9zsRH",
		Auth0URL:      "https://login.pinecone.io",
		Auth0Audience: "https://us-central1-production-console.cloudfunctions.net/api/v1",
	}

	Staging = EnvironmentConnectionSettings{
		DashboardUrl:             "https://staging.console-api.pinecone.io",
		IndexControlPlaneUrl:     "https://api-staging.pinecone.io",
		AssistantControlPlaneUrl: "https://api-staging.pinecone.io",
		AssistantDataPlaneUrl:    "https://staging-data.ke.pinecone.io",

		Auth0ClientId: "jnuhtpQxTzYw0zrpWdFUEMXS9Bx4FDAR",
		Auth0URL:      "https://internal-beta-pinecone-io.us.auth0.com",
		Auth0Audience: "https://us-central1-console-dev.cloudfunctions.net/api/v1",
	}

	DevDan = EnvironmentConnectionSettings{
		DashboardUrl:             "https://development.console-api.pinecone.io/v2",
		IndexControlPlaneUrl:     "https://api-dev.pinecone.io",
		AssistantControlPlaneUrl: "https://api-dev.pinecone.io",
		AssistantDataPlaneUrl:    "https://staging-data.ke.pinecone.io",

		Auth0ClientId: "ps7c53UDwqoLeXwUahg2A81qSMknRavi",
		Auth0URL:      "https://dev-dan-pinecone-io.us.auth0.com",
		Auth0Audience: "https://us-central1-development-pinecone.cloudfunctions.net/api/v1",
	}
)

func GetEnvConfig(env string) (EnvironmentConnectionSettings, error) {
	if env == "production" {
		return Prod, nil
	}

	if env == "staging" {
		return Staging, nil
	}

	if env == "dev-dan" {
		return DevDan, nil
	}

	return EnvironmentConnectionSettings{}, fmt.Errorf("unknown environment: %s", env)
}
