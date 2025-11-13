//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/test/e2e/helpers"
)

func TestAuthServiceAccountConfigureAndStatus(t *testing.T) {
	helpers.RequireE2E(t)
	clientID, clientSecret := helpers.RequireServiceAccount(t)

	cli := helpers.NewCLI(t)

	// Configure service account and verify org and project context are set
	ctx := context.Background()
	var context state.TargetContext
	_, err := cli.RunJSONCtx(ctx, &context,
		"auth", "configure",
		"--client-id", clientID,
		"--client-secret", clientSecret,
		"--project-id", helpers.ProjectID(),
		"--prompt-if-missing=false",
	)
	if err != nil {
		t.Fatalf("auth configure/status failed: %v", err)
	}

	if context.Organization.Id == "" || context.Project.Id == "" {
		t.Fatalf("expected TargetContext to have an Organization and Project after configuring service account credentials, got: %+v", context)
	}
}
