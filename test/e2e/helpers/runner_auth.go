//go:build e2e

package helpers

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
)

func (c *CLI) AuthConfigureServiceAccount(ctx context.Context, clientID, clientSecret, projectID string) (state.TargetContext, error) {
	var out state.TargetContext
	_, err := c.RunJSONCtx(ctx, &out,
		"auth", "configure",
		"--client-id", clientID,
		"--client-secret", clientSecret,
		"--project-id", projectID,
		"--prompt-if-missing=false",
	)
	if err != nil {
		return state.TargetContext{}, err
	}
	return out, nil
}

func (c *CLI) AuthStatus(ctx context.Context) (presenters.AuthStatus, error) {
	var st presenters.AuthStatus
	_, err := c.RunJSONCtx(ctx, &st, "auth", "status")
	return st, err
}
