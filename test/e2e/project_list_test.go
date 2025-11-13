//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/pinecone-io/cli/test/e2e/helpers"
)

func TestProjectList(t *testing.T) {
	helpers.RequireE2E(t)
	// Requires admin client
	_, _ = helpers.RequireServiceAccount(t)

	cli := helpers.NewCLI(t)

	ctx := context.Background()
	projects, err := cli.ProjectList(ctx)
	if err != nil {
		t.Fatalf("project list failed: %v", err)
	}
	if len(projects) == 0 {
		t.Fatalf("expected at least one project")
	}
}
