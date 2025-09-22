package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"golang.org/x/oauth2"
)

func PrintAuthStatus(authStatus AuthStatus) {
	log.Info().
		Str("org", authStatus.OrganizationName).
		Str("project", authStatus.ProjectName).
		Msg("Printing target context")
	writer := NewTabWriter()

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Authentication Mode\t%s\n", labelUnsetIfEmpty(authStatus.AuthMode))
	pcio.Fprintf(writer, "Global API Key\t%s\n", labelUnsetIfEmpty(authStatus.GlobalAPIKey))
	pcio.Fprintf(writer, "Service Account Client ID\t%s\n", labelUnsetIfEmpty(authStatus.ClientID))
	pcio.Fprintf(writer, "Service Account Client Secret\t%s\n", labelUnsetIfEmpty(authStatus.ClientSecret))
	pcio.Fprintf(writer, "Token Expiry\t%s\n", labelUnsetIfEmpty(authStatus.TokenExpiry))
	pcio.Fprintf(writer, "Token Time Remaining\t%s\n", labelUnsetIfEmpty(authStatus.TokenTimeRemaining))
	pcio.Fprintf(writer, "Token Scope\t%s\n", labelUnsetIfEmpty(authStatus.TokenScope))
	pcio.Fprintf(writer, "Token Organization ID\t%s\n", labelUnsetIfEmpty(authStatus.TokenOrganizationID))
	pcio.Fprintf(writer, "Environment\t%s\n", labelUnsetIfEmpty(authStatus.Environment))

	writer.Flush()
}

type AuthStatus struct {
	AuthMode            string        `json:"auth_mode,omitempty"`
	ClientID            string        `json:"client_id,omitempty"`
	ClientSecret        string        `json:"client_secret,omitempty"`
	Environment         string        `json:"environment,omitempty"`
	GlobalAPIKey        string        `json:"global_api_key,omitempty"`
	OrganizationName    string        `json:"organization_name,omitempty"`
	ProjectName         string        `json:"project_name,omitempty"`
	Token               *oauth2.Token `json:"token,omitempty"`
	TokenExpiry         string        `json:"token_expiry,omitempty"`
	TokenOrganizationID string        `json:"token_organization_id,omitempty"`
	TokenScope          string        `json:"token_scope,omitempty"`
	TokenTimeRemaining  string        `json:"token_time_remaining,omitempty"`
}
