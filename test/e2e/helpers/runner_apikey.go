//go:build e2e

package helpers

import (
	"context"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func (c *CLI) APIKeyCreate(ctx context.Context, projectID, name string) (pinecone.APIKeyWithSecret, error) {
	var created pinecone.APIKeyWithSecret
	_, err := c.RunJSONCtx(ctx, &created, "api-key", "create", "--id", projectID, "--name", name)
	return created, err
}

func (c *CLI) APIKeyDescribe(ctx context.Context, id string) (pinecone.APIKey, error) {
	var d pinecone.APIKey
	_, err := c.RunJSONCtx(ctx, &d, "api-key", "describe", "--id", id)
	return d, err
}

func (c *CLI) APIKeyDelete(ctx context.Context, id string) error {
	_, _, err := c.RunCtx(ctx, "api-key", "delete", "--id", id, "--skip-confirmation")
	return err
}
