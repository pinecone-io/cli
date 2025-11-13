//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/pinecone-io/cli/test/e2e/helpers"
)

func TestAPIKeyLifecycle(t *testing.T) {
	helpers.RequireE2E(t)
	// Requires admin client to manage API keys
	_, _ = helpers.RequireServiceAccount(t)

	projID := helpers.ProjectID()
	if projID == "" {
		t.Skip("PC_E2E_PROJECT_ID not set; skipping api-key lifecycle test")
	}

	cli := helpers.NewCLI(t)
	name := helpers.RandomName("e2e-key")
	ctx := context.Background()

	// Create
	created, err := cli.APIKeyCreate(ctx, projID, name)
	if err != nil {
		t.Fatalf("api-key create failed: %v", err)
	}
	if created.Key.Id == "" || created.Value == "" {
		t.Fatalf("expected created key with secret, got: %+v", created)
	}

	// Describe
	desc, err := cli.APIKeyDescribe(ctx, created.Key.Id)
	if err != nil {
		t.Fatalf("api-key describe failed: %v", err)
	}
	if desc.Id != created.Key.Id {
		t.Fatalf("describe id mismatch: expected %s got %s", created.Key.Id, desc.Id)
	}

	// Delete / Ensure cleanup
	t.Cleanup(func() {
		_ = cli.APIKeyDelete(ctx, created.Key.Id)
	})
}
