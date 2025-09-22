package presenters

import (
	"strings"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/auth"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"golang.org/x/oauth2"
)

func PrintAuthStatus(token *oauth2.Token) {
	authMode := string(state.TargetCreds.Get().AuthContext)
	orgName := state.TargetOrg.Get().Name
	projName := state.TargetProj.Get().Name
	environment := config.Environment.Get()

	// Global API Key
	globalAPIKey := secrets.GlobalApiKey.Get()

	// Service Account
	clientId := secrets.ClientId.Get()
	clientSecret := secrets.ClientSecret.Get()

	// Extract token information
	var claims *auth.MyCustomClaims
	expStr := ""
	remaining := ""
	scope := ""
	orgId := ""

	if token != nil {
		if token.AccessToken != "" {
			claims, _ = auth.ParseClaimsUnverified(token)
		}

		if !token.Expiry.IsZero() {
			expStr = token.Expiry.Format(time.RFC3339)
			remaining = time.Until(token.Expiry).Round(time.Second).String()
		}

		if claims != nil {
			scope = claims.Scope
			orgId = claims.OrgId
		}
	}

	log.Info().
		Str("org", orgName).
		Str("project", projName).
		Msg("Printing target context")
	writer := NewTabWriter()

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Authentication Mode\t%s\n", labelUnsetIfEmpty(string(authMode)))
	pcio.Fprintf(writer, "Global API Key\t%s\n", labelUnsetIfEmpty(globalAPIKey))
	pcio.Fprintf(writer, "Service Account Client ID\t%s\n", labelUnsetIfEmpty(clientId))
	pcio.Fprintf(writer, "Service Account Client Secret\t%s\n", labelUnsetIfEmpty(clientSecret))
	pcio.Fprintf(writer, "Token Expiry\t%s\n", labelUnsetIfEmpty(expStr))
	pcio.Fprintf(writer, "Token Time Remaining\t%s\n", labelUnsetIfEmpty(remaining))
	pcio.Fprintf(writer, "Token Scope\t%s\n", labelUnsetIfEmpty(scope))
	pcio.Fprintf(writer, "Token Organization ID\t%s\n", labelUnsetIfEmpty(orgId))
	pcio.Fprintf(writer, "Environment\t%s\n", labelUnsetIfEmpty(environment))

	writer.Flush()
}
