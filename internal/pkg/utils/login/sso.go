package login

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/environment"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
)

// dashboardOrg is the subset of the dashboard API org response needed for SSO lookup.
type dashboardOrg struct {
	Id                string `json:"id"`
	SSOConnectionName string `json:"sso_connection_name"`
	EnforceSSO        bool   `json:"enforce_sso_authentication"`
}

type dashboardOrgsResponse struct {
	NewOrgs []dashboardOrg `json:"newOrgs"`
}

// FetchSSOConnection calls the private dashboard API to retrieve the Auth0
// connection name for the given orgId. It returns ("", nil) when the org has
// no SSO configured, enforce_sso_authentication is false, or any error occurs.
// Errors are non-fatal: the caller should proceed with a normal login URL.
func FetchSSOConnection(ctx context.Context, orgId string) (string, error) {
	token, err := oauth.Token(ctx)
	if err != nil || token == nil || token.AccessToken == "" {
		log.Debug().Str("orgId", orgId).Msg("SSO lookup skipped: no valid token available")
		return "", nil
	}

	envConfig, err := environment.GetEnvConfig(config.Environment.Get())
	if err != nil {
		return "", nil
	}

	return fetchSSOConnectionFromURL(ctx, orgId, token.AccessToken, http.DefaultClient, envConfig.DashboardUrl)
}

// fetchSSOConnectionFromURL is the testable core: it takes an explicit HTTP
// client and dashboard base URL so tests can inject a local httptest.Server.
func fetchSSOConnectionFromURL(ctx context.Context, orgId string, accessToken string, client *http.Client, dashboardURL string) (string, error) {
	url := dashboardURL + "/v2/dashboard/organizations"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", nil
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Str("orgId", orgId).Msg("SSO lookup: dashboard API request failed")
		return "", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Debug().Int("status", resp.StatusCode).Str("orgId", orgId).Msg("SSO lookup: dashboard API returned non-2xx")
		return "", nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		log.Debug().Err(err).Str("orgId", orgId).Msg("SSO lookup: failed to read dashboard API response")
		return "", nil
	}

	var orgsResp dashboardOrgsResponse
	if err := json.Unmarshal(body, &orgsResp); err != nil {
		log.Debug().Err(err).Str("orgId", orgId).Msg("SSO lookup: failed to decode dashboard API response")
		return "", nil
	}

	for _, org := range orgsResp.NewOrgs {
		if org.Id == orgId {
			if org.EnforceSSO && org.SSOConnectionName != "" {
				log.Debug().Str("orgId", orgId).Str("connection", org.SSOConnectionName).Msg("SSO lookup: found connection")
				return org.SSOConnectionName, nil
			}
			return "", nil
		}
	}

	log.Debug().Str("orgId", orgId).Msg("SSO lookup: org not found in dashboard response")
	return "", nil
}
