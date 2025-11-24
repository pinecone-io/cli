package inputpolicy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath_AllowsRegularJSON(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "ok.json")
	if err := os.WriteFile(fp, []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePath(fp); err != nil {
		t.Fatalf("expected ok, got %v", err)
	}
}

func TestValidatePath_RejectsDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := ValidatePath(dir); err == nil {
		t.Fatalf("expected error for directory")
	}
}
