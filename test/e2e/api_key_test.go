//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/pinecone-io/cli/test/e2e/helpers"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
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
	var create pinecone.APIKeyWithSecret
	_, err := cli.RunJSONCtx(ctx, &create, "api-key", "create", "--id", projID, "--name", name)
	if err != nil {
		t.Fatalf("api-key create failed: %v", err)
	}
	if create.Key.Id == "" || create.Value == "" {
		t.Fatalf("expected created key with secret, got: %+v", create)
	}

	// Describe
	var desc pinecone.APIKey
	_, err = cli.RunJSONCtx(ctx, &desc, "api-key", "describe", "--id", create.Key.Id)
	if err != nil {
		t.Fatalf("api-key describe failed: %v", err)
	}
	if desc.Id != create.Key.Id {
		t.Fatalf("describe id mismatch: expected %s got %s", create.Key.Id, desc.Id)
	}

	// Delete / Ensure cleanup
	t.Cleanup(func() {
		_, _, err := cli.RunCtx(ctx, "api-key", "delete", "--id", desc.Id, "--skip-confirmation")
		if err != nil {
			t.Fatalf("api-key delete failed: %v", err)
		}
	})
}
