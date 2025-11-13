//go:build e2e

package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/pinecone-io/cli/test/e2e/helpers"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func TestIndexServerless_ServiceAccount(t *testing.T) {
	helpers.RequireE2E(t)
	_, _ = helpers.RequireServiceAccount(t)
	runIndexServerlessLifecycle(t)
}

func TestIndexServerless_APIKey(t *testing.T) {
	helpers.RequireE2E(t)
	_ = helpers.RequireAPIKey(t)
	runIndexServerlessLifecycle(t)
}

func runIndexServerlessLifecycle(t *testing.T) {
	cli := helpers.NewCLI(t)
	ctx := context.Background()
	name := helpers.RandomName("e2e-srvless")

	// Create serverless index
	created, err := cli.IndexCreateServerless(ctx, helpers.IndexCreateServerlessOptions{
		Name:      name,
		Cloud:     helpers.Cloud(),
		Region:    helpers.Region(),
		Dimension: helpers.Dimension(),
	})
	if err != nil {
		t.Fatalf("index create failed: %v", err)
	}
	if created.Name != name {
		t.Fatalf("created index name mismatch: expected %s got %s", name, created.Name)
	}

	// Ensure cleanup
	t.Cleanup(func() {
		_ = cli.IndexDelete(ctx, name)
	})

	// Wait for readiness
	if err := helpers.WaitForIndexReady(cli, name, 5*time.Minute); err != nil {
		t.Fatalf("index not ready: %v", err)
	}

	// Describe
	desc, err := cli.IndexDescribe(ctx, name)
	if err != nil {
		t.Fatalf("index describe failed: %v", err)
	}
	if desc.Name != name {
		t.Fatalf("describe name mismatch: expected %s got %s", name, desc.Name)
	}

	// List and assert presence (len > 0 and name contained)
	list, stdout, err := helpers.MustRunJSON[[]pinecone.Index](cli, ctx, "index", "list")
	if err != nil {
		t.Fatalf("index list failed: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("expected at least one index in list")
	}
	if !strings.Contains(stdout, "\"name\": \""+name+"\"") {
		t.Fatalf("created index not found in list output")
	}

	// Delete explicitly (also registered in cleanup)
	err = cli.IndexDelete(ctx, name)
	if err != nil {
		t.Fatalf("index delete failed: %v", err)
	}
}
