package argio

import (
	"io"
	"strings"
	"testing"
)

func TestOpenArgReader_Inline(t *testing.T) {
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

func TestReadAll_Inline(t *testing.T) {
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
