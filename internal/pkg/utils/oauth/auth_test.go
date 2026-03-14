package oauth

import (
	"context"
	"net/url"
	"testing"
)

func TestSourceTagConstant(t *testing.T) {
	if SourceTag != "pinecone_cli" {
		t.Errorf("expected SourceTag to be %q, got %q", "pinecone_cli", SourceTag)
	}
}

func TestGetAuthURL_ContainsSourceTag(t *testing.T) {
	a := &Auth{}
	ctx := context.Background()

	_, challenge, err := a.CreateNewVerifierAndChallenge()
	if err != nil {
		t.Fatalf("failed to create verifier/challenge: %v", err)
	}

	rawURL, err := a.GetAuthURL(ctx, "test-csrf-state", challenge, nil)
	if err != nil {
		t.Fatalf("GetAuthURL returned error: %v", err)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("failed to parse auth URL: %v", err)
	}

	params := parsed.Query()

	if got := params.Get("sourceTag"); got != SourceTag {
		t.Errorf("expected sourceTag=%q, got %q", SourceTag, got)
	}
}

func TestGetAuthURL_RequiredParams(t *testing.T) {
	a := &Auth{}
	ctx := context.Background()

	_, challenge, err := a.CreateNewVerifierAndChallenge()
	if err != nil {
		t.Fatalf("failed to create verifier/challenge: %v", err)
	}

	csrfState := "test-state-123"
	rawURL, err := a.GetAuthURL(ctx, csrfState, challenge, nil)
	if err != nil {
		t.Fatalf("GetAuthURL returned error: %v", err)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("failed to parse auth URL: %v", err)
	}

	params := parsed.Query()

	if params.Get("state") != csrfState {
		t.Errorf("expected state=%q, got %q", csrfState, params.Get("state"))
	}
	if params.Get("code_challenge") != challenge {
		t.Errorf("expected code_challenge=%q, got %q", challenge, params.Get("code_challenge"))
	}
	if params.Get("code_challenge_method") != "S256" {
		t.Errorf("expected code_challenge_method=S256, got %q", params.Get("code_challenge_method"))
	}
	if params.Get("audience") == "" {
		t.Error("expected audience param to be set")
	}
}

func TestGetAuthURL_WithOrgId(t *testing.T) {
	a := &Auth{}
	ctx := context.Background()

	_, challenge, err := a.CreateNewVerifierAndChallenge()
	if err != nil {
		t.Fatalf("failed to create verifier/challenge: %v", err)
	}

	orgId := "test-org-456"
	rawURL, err := a.GetAuthURL(ctx, "state", challenge, &orgId)
	if err != nil {
		t.Fatalf("GetAuthURL returned error: %v", err)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("failed to parse auth URL: %v", err)
	}

	if got := parsed.Query().Get("orgId"); got != orgId {
		t.Errorf("expected orgId=%q, got %q", orgId, got)
	}
}

func TestGetAuthURL_WithEmptyOrgId(t *testing.T) {
	a := &Auth{}
	ctx := context.Background()

	_, challenge, err := a.CreateNewVerifierAndChallenge()
	if err != nil {
		t.Fatalf("failed to create verifier/challenge: %v", err)
	}

	emptyOrgId := ""
	rawURL, err := a.GetAuthURL(ctx, "state", challenge, &emptyOrgId)
	if err != nil {
		t.Fatalf("GetAuthURL returned error: %v", err)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("failed to parse auth URL: %v", err)
	}

	if got := parsed.Query().Get("orgId"); got != "" {
		t.Errorf("expected orgId to be absent for empty string, got %q", got)
	}
}
