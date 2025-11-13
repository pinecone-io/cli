//go:build e2e

package e2e

import (
	"context"
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
	args := []string{
		"index", "create",
		"--name", name,
		"--cloud", helpers.Cloud(),
		"--region", helpers.Region(),
		"--dimension", strconvI(helpers.Dimension()),
		"--metric", "cosine",
	}
	var idx pinecone.Index
	_, err := cli.RunJSONCtx(ctx, &idx, args...)
	if err != nil {
		t.Fatalf("index create failed: %v", err)
	}
	if idx.Name != name {
		t.Fatalf("created index name mismatch: expected %s got %s", name, idx.Name)
	}

	// Wait for readiness
	if err := helpers.WaitForIndexReady(cli, name, 5*time.Minute); err != nil {
		t.Fatalf("index not ready: %v", err)
	}

	// Describe
	var desc pinecone.Index
	_, err = cli.RunJSONCtx(ctx, &desc, "index", "describe", "--name", name)
	if err != nil {
		t.Fatalf("index describe failed: %v", err)
	}
	if desc.Name != name {
		t.Fatalf("describe name mismatch: expected %s got %s", name, desc.Name)
	}

	// List and assert presence (len > 0 and name contained)
	var list []pinecone.Index
	_, err = cli.RunJSONCtx(ctx, &list, "index", "list")
	if err != nil {
		t.Fatalf("index list failed: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("expected at least one index in list")
	}

	newIdxListed := false
	for _, idx := range list {
		if idx.Name == name {
			newIdxListed = true
		}
	}
	if !newIdxListed {
		t.Fatalf("created index not found in list output")
	}

	// Delete / Ensure cleanup
	t.Cleanup(func() {
		_, _, err = cli.RunCtx(ctx, "index", "delete", "--name", name)
		if err != nil {
			t.Fatalf("index delete failed: %v", err)
		}
	})
}

// local helper (duplicated here to avoid importing strconv in tests)
func strconvI(n int) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = digits[n%10]
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
