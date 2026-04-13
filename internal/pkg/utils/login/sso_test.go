package login

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newDashboardServer starts an httptest.Server that returns the given org list.
// Pass a non-zero statusCode to simulate an error response.
func newDashboardServer(t *testing.T, orgs []dashboardOrg, statusCode int) *httptest.Server {
	t.Helper()
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if statusCode != http.StatusOK {
			http.Error(w, "error", statusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(dashboardOrgsResponse{NewOrgs: orgs})
	}))
}

func TestFetchSSOConnection_EnforcedWithConnection(t *testing.T) {
	server := newDashboardServer(t, []dashboardOrg{
		{Id: "org-1", SSOConnectionName: "alby-saml", EnforceSSO: true},
	}, 0)
	defer server.Close()

	conn, err := fetchSSOConnectionFromURL(context.Background(), "org-1", "fake-token", server.Client(), server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn != "alby-saml" {
		t.Errorf("expected %q, got %q", "alby-saml", conn)
	}
}

func TestFetchSSOConnection_NotEnforced(t *testing.T) {
	server := newDashboardServer(t, []dashboardOrg{
		{Id: "org-1", SSOConnectionName: "alby-saml", EnforceSSO: false},
	}, 0)
	defer server.Close()

	conn, err := fetchSSOConnectionFromURL(context.Background(), "org-1", "fake-token", server.Client(), server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn != "" {
		t.Errorf("expected empty connection when SSO not enforced, got %q", conn)
	}
}

func TestFetchSSOConnection_OrgNotFound(t *testing.T) {
	server := newDashboardServer(t, []dashboardOrg{
		{Id: "org-other", SSOConnectionName: "other-saml", EnforceSSO: true},
	}, 0)
	defer server.Close()

	conn, err := fetchSSOConnectionFromURL(context.Background(), "org-1", "fake-token", server.Client(), server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn != "" {
		t.Errorf("expected empty connection when org not found, got %q", conn)
	}
}

func TestFetchSSOConnection_NonOKStatus(t *testing.T) {
	server := newDashboardServer(t, nil, http.StatusUnauthorized)
	defer server.Close()

	conn, err := fetchSSOConnectionFromURL(context.Background(), "org-1", "fake-token", server.Client(), server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn != "" {
		t.Errorf("expected empty connection on non-2xx response, got %q", conn)
	}
}

func TestFetchSSOConnection_EmptyConnectionName(t *testing.T) {
	server := newDashboardServer(t, []dashboardOrg{
		{Id: "org-1", SSOConnectionName: "", EnforceSSO: true},
	}, 0)
	defer server.Close()

	conn, err := fetchSSOConnectionFromURL(context.Background(), "org-1", "fake-token", server.Client(), server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn != "" {
		t.Errorf("expected empty connection when name is empty, got %q", conn)
	}
}

func TestFetchSSOConnection_MultipleOrgs(t *testing.T) {
	server := newDashboardServer(t, []dashboardOrg{
		{Id: "org-1", SSOConnectionName: "org1-saml", EnforceSSO: true},
		{Id: "org-2", SSOConnectionName: "org2-saml", EnforceSSO: true},
	}, 0)
	defer server.Close()

	conn, err := fetchSSOConnectionFromURL(context.Background(), "org-2", "fake-token", server.Client(), server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn != "org2-saml" {
		t.Errorf("expected %q, got %q", "org2-saml", conn)
	}
}
