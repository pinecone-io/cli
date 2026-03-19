package testutils

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// CaptureStdout redirects os.Stdout to a pipe for the duration of f,
// returning everything written to stdout as a trimmed string.
func CaptureStdout(t *testing.T, f func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("CaptureStdout: os.Pipe: %v", err)
	}

	prev := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = prev }()

	f()
	w.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("CaptureStdout: reading pipe: %v", err)
	}
	return strings.TrimSpace(buf.String())
}
