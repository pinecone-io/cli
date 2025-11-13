//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/test/e2e/helpers"
)

func TestTargetSetAndShow(t *testing.T) {
	helpers.RequireE2E(t)
	// admin needed to resolve org/project by ID
	_, _ = helpers.RequireServiceAccount(t)

	orgID := helpers.OrgID()
	projID := helpers.ProjectID()
	if orgID == "" || projID == "" {
		t.Skip("PC_E2E_ORG_ID or PC_E2E_PROJECT_ID not set; skipping target test")
	}

	cli := helpers.NewCLI(t)

	// Set target
	ctx := context.Background()
	_, _, err := cli.RunCtx(ctx, "target", "--org-id", orgID, "--project-id", projID)
	if err != nil {
		t.Fatalf("target set failed: %v", err)
	}

	// Show target and verify
	var tc state.TargetContext
	_, err = cli.RunJSONCtx(ctx, &tc, "target", "--show")
	if err != nil {
		t.Fatalf("target --show failed: %v", err)
	}
	if tc.Organization.Id != orgID {
		t.Fatalf("expected org id %s, got %s", orgID, tc.Organization.Id)
	}
	if tc.Project.Id != projID {
		t.Fatalf("expected project id %s, got %s", projID, tc.Project.Id)
	}
}
