//go:build e2e

package helpers

import (
	"context"
	"testing"
)

// WithServerlessIndex creates a serverless index with a randomized name and
// ensures cleanup after the test via t.Cleanup.
func WithServerlessIndex(t *testing.T, cli *CLI, opts IndexCreateServerlessOptions, fn func(name string)) {
	t.Helper()
	ctx := context.Background()
	name := RandomName("e2e-srvless")
	opts.Name = name

	_, err := cli.IndexCreateServerless(ctx, opts)
	if err != nil {
		t.Fatalf("failed to create serverless index: %v", err)
	}
	t.Cleanup(func() {
		_ = cli.IndexDelete(ctx, name)
	})
	fn(name)
}
