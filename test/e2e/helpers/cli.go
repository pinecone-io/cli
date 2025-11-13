//go:build e2e

package helpers

import (
	"context"
	"encoding/json"

	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// CLI is a wrapper around the CLI binary and environment variables.
// It provides methods for running CLI commands and parsing various kinds of output in tests.
type CLI struct {
	// Bin is the path to the CLI binary which will be used to run commands.
	Bin string
	// BaseEnv are the environment variables that will be used to run commands.
	BaseEnv []string
	// T is the testing.T instance that will be used to log messages.
	T       *testing.T
	Timeout time.Duration
	Debug   bool
}

// NewCLI resolves the CLI binary and returns a runner that can be used to exercise commands.
// Binary Resolution order:
// 1) PC_BIN if set (validated)
// 2) (Optional) PC_E2E_USE_PATH=1 to use 'pc' from local PATH
// 3) Otherwise fail (TestMain or caller should provide PC_BIN)
func NewCLI(t *testing.T) *CLI {
	// If PC_BIN is set, validate it's usable.
	// Otherwise, check if PC_E2E_USE_PATH is set and check for 'pc' in the PATH.
	bin := os.Getenv("PC_BIN")
	if bin != "" {
		if _, err := os.Stat(bin); err != nil {
			t.Fatalf("PC_BIN is not usable: %s: %v", bin, err)
		}
	} else if os.Getenv("PC_E2E_USE_PATH") == "1" {
		p, err := exec.LookPath("pc")
		if err != nil {
			t.Fatalf("PC_E2E_USE_PATH=1 is set, but 'pc' not found in PATH: %v", err)
		}
		bin = p
	}

	if bin == "" {
		t.Fatalf("pc binary not found. Set PC_BIN or PC_E2E_USE_PATH=1")
	}

	// Copy all environment variables into baseEnv, overriding PC_E2E_HOME if set.
	baseEnv := os.Environ()
	if os.Getenv("PC_E2E_HOME") != "" {
		baseEnv = appendOrReplaceEnv(baseEnv, "PC_E2E_HOME", os.Getenv("PC_E2E_HOME"))
	}

	return &CLI{
		Bin:     bin,
		BaseEnv: baseEnv,
		T:       t,
		Timeout: 2 * time.Minute,
		Debug:   os.Getenv("PC_E2E_DEBUG") == "1",
	}
}

func (c *CLI) Run(args ...string) (string, string, error) {
	c.T.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	return c.RunCtx(ctx, args...)
}

// RunJSON appends --json if not present and unmarshals stdout into out.
func (c *CLI) RunJSON(out any, args ...string) (string, error) {
	c.T.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	return c.RunJSONCtx(ctx, out, args...)
}

// RunCtx executes the CLI with a provided context.
func (c *CLI) RunCtx(ctx context.Context, args ...string) (string, string, error) {
	c.T.Helper()
	cmd := exec.CommandContext(ctx, c.Bin, args...)
	cmd.Env = c.BaseEnv
	if c.Debug {
		c.T.Logf("RUN: %s %s", c.Bin, strings.Join(args, " "))
	}
	out, err := cmd.CombinedOutput()
	stdout := string(out)

	// Cobra writes most output to stdout via their pcio helpers; retain single stream
	// We still return stderr string for API compatibility.
	stderr := ""
	if ctx.Err() == context.DeadlineExceeded {
		if c.Debug {
			c.T.Logf("TIMEOUT stdout/stderr:\n%s", stdout)
		}
		return stdout, stderr, errors.New("command timed out")
	}
	if err != nil {
		// Always log output on error when debug is enabled
		if c.Debug {
			c.T.Logf("ERR stdout/stderr:\n%s", stdout)
		}
		return stdout, stderr, fmt.Errorf("run error: %w\nstdout/stderr:\n%s", err, stdout)
	}
	return stdout, stderr, nil
}

// RunJSONCtx appends --json if not present and unmarshals stdout into out, honoring context.
func (c *CLI) RunJSONCtx(ctx context.Context, out any, args ...string) (string, error) {
	c.T.Helper()
	if !hasJSONFlag(args) {
		args = append(args, "--json")
	}
	stdout, _, err := c.RunCtx(ctx, args...)
	if err != nil {
		return stdout, err
	}
	if err := json.Unmarshal([]byte(stdout), out); err != nil {
		return stdout, fmt.Errorf("failed to parse JSON: %w\nraw: %s", err, stdout)
	}
	return stdout, nil
}

// MustRunJSON executes and decodes JSON into the provided generic type.
// func MustRunJSON[T any](c *CLI, ctx context.Context, args ...string) (T, string, error) {
// 	var zero T
// 	var v T
// 	stdout, err := c.RunJSONCtx(ctx, &v, args...)
// 	if err != nil {
// 		return zero, stdout, err
// 	}
// 	return v, stdout, nil
// }

func hasJSONFlag(args []string) bool {
	for _, a := range args {
		if a == "--json" {
			return true
		}
	}
	return false
}

func appendOrReplaceEnv(env []string, key, value string) []string {
	prefix := key + "="
	for i, e := range env {
		if strings.HasPrefix(e, prefix) {
			env[i] = prefix + value
			return env
		}
	}
	return append(env, prefix+value)
}
