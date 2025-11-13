//go:build e2e

package helpers

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
)

func (c *CLI) TargetSetByIDs(ctx context.Context, orgID, projectID string) error {
	_, _, err := c.RunCtx(ctx, "target", "--org-id", orgID, "--project-id", projectID)
	return err
}

func (c *CLI) TargetShow(ctx context.Context) (state.TargetContext, error) {
	var tc state.TargetContext
	_, err := c.RunJSONCtx(ctx, &tc, "target", "--show")
	return tc, err
}
