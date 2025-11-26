package presenters

import (
	"bytes"
	"strings"
	"testing"
	"text/tabwriter"
)

func TestPrintEmptyState(t *testing.T) {
	var buf bytes.Buffer
	writer := tabwriter.NewWriter(&buf, 12, 1, 4, ' ', 0)

	if !PrintEmptyState(writer, "test data") {
		t.Fatalf("PrintEmptyState should always return true")
	}

	got := strings.TrimSpace(buf.String())
	if got != "No test data available." {
		t.Fatalf("unexpected output: %q", got)
	}
}
