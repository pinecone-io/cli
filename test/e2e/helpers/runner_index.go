//go:build e2e

package helpers

import (
	"context"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

type IndexCreateServerlessOptions struct {
	Name      string
	Cloud     string
	Region    string
	Dimension int
	Metric    string // optional
}

func (c *CLI) IndexCreateServerless(ctx context.Context, opts IndexCreateServerlessOptions) (pinecone.Index, error) {
	args := []string{
		"index", "create",
		"--name", opts.Name,
		"--cloud", opts.Cloud,
		"--region", opts.Region,
		"--dimension", strconvI(opts.Dimension),
	}
	if opts.Metric != "" {
		args = append(args, "--metric", opts.Metric)
	}
	var idx pinecone.Index
	_, err := c.RunJSONCtx(ctx, &idx, args...)
	return idx, err
}

func (c *CLI) IndexDescribe(ctx context.Context, name string) (pinecone.Index, error) {
	var idx pinecone.Index
	_, err := c.RunJSONCtx(ctx, &idx, "index", "describe", "--name", name)
	return idx, err
}

func (c *CLI) IndexList(ctx context.Context) ([]pinecone.Index, error) {
	var idxs []pinecone.Index
	_, err := c.RunJSONCtx(ctx, &idxs, "index", "list")
	return idxs, err
}

func (c *CLI) IndexDelete(ctx context.Context, name string) error {
	_, _, err := c.RunCtx(ctx, "index", "delete", "--name", name)
	return err
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
