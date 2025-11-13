//go:build e2e

package helpers

import (
	"context"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func (c *CLI) ProjectList(ctx context.Context) ([]pinecone.Project, error) {
	var projects []pinecone.Project
	_, err := c.RunJSONCtx(ctx, &projects, "project", "list")
	return projects, err
}
