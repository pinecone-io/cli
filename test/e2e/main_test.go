//go:build e2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// e2e Testing entry point. TestMain performs per-package setup. It builds the CLI binary once if PC_BIN
// is not provided and PC_E2E_USE_PATH is not set, then sets PC_BIN for tests.
func TestMain(m *testing.M) {
	// If an explicit binary is not provided, and we're not using what's in the PATH, build the CLI binary
	if os.Getenv("PC_BIN") == "" && os.Getenv("PC_E2E_USE_PATH") != "1" {
		root := findRepoRoot()
		out := filepath.Join(root, "dist", "pc-e2e")
		_ = os.MkdirAll(filepath.Dir(out), 0755)

		cmd := exec.Command("go", "build", "-o", out, "./cmd/pc")
		cmd.Dir = root
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
		_ = os.Setenv("PC_BIN", out)
	}

	// Set service account credentials and target context
	if os.Getenv("PINECONE_CLIENT_ID") != "" && os.Getenv("PINECONE_CLIENT_SECRET") != "" {
		args := []string{
			os.Getenv("PC_BIN"), "auth", "configure",
			"--client-id", os.Getenv("PINECONE_CLIENT_ID"),
			"--client-secret", os.Getenv("PINECONE_CLIENT_SECRET"),
			"--prompt-if-missing=false",
			"--json",
		}
		if os.Getenv("PC_E2E_PROJECT_ID") != "" {
			args = append(args, "--project-id", os.Getenv("PC_E2E_PROJECT_ID"))
		}
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	// Isolate $HOME so we don't blow away local ~/.config/pinecone.
	// Setting $HOME to a temporary directory will isolate Viper configurations
	tempHome, _ := os.MkdirTemp("", "pc-e2e-home-*")
	_ = os.Setenv("PC_E2E_HOME", tempHome)

	// Run tests
	code := m.Run()

	_ = os.RemoveAll(tempHome)
	os.Exit(code)
}

// findRepoRoot walks up from the current working directory to find go.mod.
func findRepoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("go.mod not found")
		}
		dir = parent
	}
}
