//go:build e2e

package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strconv"
	"testing"
	"time"
)

// E2EConfig captures suite parameters parsed once from environment.
type E2EConfig struct {
	Cloud     string
	Region    string
	Dimension int
	OrgID     string
	ProjectID string
}

func ParseE2EConfig() E2EConfig {
	return E2EConfig{
		Cloud:     Cloud(),
		Region:    Region(),
		Dimension: Dimension(),
		OrgID:     OrgID(),
		ProjectID: ProjectID(),
	}
}

// PC_E2E=1 must explicitly be set to run E2E tests
func RequireE2E(t *testing.T) {
	t.Helper()
	if os.Getenv("PC_E2E") != "1" {
		t.Skip("PC_E2E != 1; skipping e2e tests")
	}
}

// Tests which require service account credentials (PINECONE_CLIENT_ID & PINECONE_CLIENT_SECRET)
func RequireServiceAccount(t *testing.T) (string, string) {
	t.Helper()
	id := os.Getenv("PINECONE_CLIENT_ID")
	secret := os.Getenv("PINECONE_CLIENT_SECRET")
	if id == "" || secret == "" {
		t.Skip("PINECONE_CLIENT_ID or PINECONE_CLIENT_SECRET not set; skipping service account tests")
	}
	return id, secret
}

// Tests which require an API key (PINECONE_API_KEY)
func RequireAPIKey(t *testing.T) string {
	t.Helper()
	k := os.Getenv("PINECONE_API_KEY")
	if k == "" {
		t.Skip("PINECONE_API_KEY not set; skipping API key tests")
	}
	return k
}

// PC_E2E_CLOUD (default: "aws")
func Cloud() string {
	if v := os.Getenv("PC_E2E_CLOUD"); v != "" {
		return v
	}
	return "aws"
}

// PC_E2E_REGION (default: "us-east-1")
func Region() string {
	if v := os.Getenv("PC_E2E_REGION"); v != "" {
		return v
	}
	return "us-east-1"
}

// PC_E2E_DIMENSION (default: 8)
func Dimension() int {
	if v := os.Getenv("PC_E2E_DIMENSION"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return 8
}

// PC_E2E_PROJECT_ID
func ProjectID() string {
	return os.Getenv("PC_E2E_PROJECT_ID")
}

// PC_E2E_ORG_ID
func OrgID() string {
	return os.Getenv("PC_E2E_ORG_ID")
}

func RandomName(prefix string) string {
	var b [4]byte
	_, _ = rand.Read(b[:])
	return prefix + "-" + time.Now().UTC().Format("20060102-150405") + "-" + hex.EncodeToString(b[:])
}
