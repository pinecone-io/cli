package argio

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/utils/stdin"
)

func Test_OpenReader_Inline(t *testing.T) {
	rc, src, err := OpenReader(`{"a":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Kind != SourceInline || src.Label != "inline" {
		t.Fatalf("unexpected source: %+v", src)
	}
	defer rc.Close()
	b, _ := io.ReadAll(rc)
	if !strings.Contains(string(b), `"a":1`) {
		t.Fatalf("unexpected body %q", string(b))
	}
}

func Test_ReadAll_Inline(t *testing.T) {
	b, src, err := ReadAll(`{"a":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Kind != SourceInline || src.Label != "inline" {
		t.Fatalf("unexpected source: %+v", src)
	}
	if !strings.Contains(string(b), `"a":1`) {
		t.Fatalf("unexpected body %q", string(b))
	}
}

func Test_ReadAll_JSONFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "payload.json")
	if err := os.WriteFile(path, []byte(`{"file":true}`), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	b, src, err := ReadAll(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Kind != SourceFile || src.Label != path {
		t.Fatalf("expected SourceFile for %s, got %+v", path, src)
	}
	if !strings.Contains(string(b), `"file":true`) {
		t.Fatalf("unexpected body %q", string(b))
	}
}

func Test_ReadAll_JSONLFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "payload.jsonl")
	content := "{\"line\":1}\n{\"line\":2}\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	rc, src, err := OpenReader(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Kind != SourceFile || src.Label != path {
		t.Fatalf("expected SourceFile for %s, got %+v", path, src)
	}
	defer rc.Close()
	data, _ := io.ReadAll(rc)
	if string(data) != content {
		t.Fatalf("unexpected data %q", string(data))
	}
}

func Test_OpenReader_DoesNotTreatInlineJSONAsFile(t *testing.T) {
	value := `{"endsWith":".json"}`
	rc, src, err := OpenReader(value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Kind != SourceInline {
		t.Fatalf("expected inline source, got %+v", src)
	}
	defer rc.Close()
	data, _ := io.ReadAll(rc)
	if !strings.Contains(string(data), `.json`) {
		t.Fatalf("unexpected data %q", string(data))
	}
}

func Test_OpenReader_StdinDash(t *testing.T) {
	stdin.ResetForTests()
	orig := os.Stdin
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdin = pr
	defer func() {
		os.Stdin = orig
		pr.Close()
		stdin.ResetForTests()
	}()

	go func() {
		defer pw.Close()
		if _, err := pw.Write([]byte(`{"stdin":true}`)); err != nil {
			panic(err)
		}
	}()

	rc, src, err := OpenReader("-")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Kind != SourceStdin {
		t.Fatalf("expected SourceStdin, got %+v", src)
	}
	data, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read stdin: %v", err)
	}
	if string(data) != `{"stdin":true}` {
		t.Fatalf("unexpected data %q", string(data))
	}
}
