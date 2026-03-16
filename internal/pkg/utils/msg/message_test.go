package msg

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

// captureStderr redirects os.Stderr to a pipe for the duration of f,
// returning everything written to stderr as a string.
func captureStderr(t *testing.T, f func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	prev := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = prev }()

	f()
	w.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("captureStderr: reading pipe: %v", err)
	}
	return buf.String()
}

func TestMsgFunctions_WriteToStderr(t *testing.T) {
	// Disable color so prefixes are plain ASCII and assertions are stable.
	prev := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = prev }()

	tests := []struct {
		name           string
		fn             func()
		expectedPrefix string
		expectedText   string
	}{
		{
			name:           "FailMsg writes [ERROR] prefix to stderr",
			fn:             func() { FailMsg("something went wrong") },
			expectedPrefix: "[ERROR]",
			expectedText:   "something went wrong",
		},
		{
			name:           "SuccessMsg writes [SUCCESS] prefix to stderr",
			fn:             func() { SuccessMsg("operation complete") },
			expectedPrefix: "[SUCCESS]",
			expectedText:   "operation complete",
		},
		{
			name:           "WarnMsg writes [WARN] prefix to stderr",
			fn:             func() { WarnMsg("take care") },
			expectedPrefix: "[WARN]",
			expectedText:   "take care",
		},
		{
			name:           "InfoMsg writes [INFO] prefix to stderr",
			fn:             func() { InfoMsg("just so you know") },
			expectedPrefix: "[INFO]",
			expectedText:   "just so you know",
		},
		{
			name:           "HintMsg writes Hint: prefix to stderr",
			fn:             func() { HintMsg("try this") },
			expectedPrefix: "Hint:",
			expectedText:   "try this",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := captureStderr(t, tt.fn)
			assert.Contains(t, out, tt.expectedPrefix)
			assert.Contains(t, out, tt.expectedText)
		})
	}
}

func TestMsgFunctions_FormatString(t *testing.T) {
	// Disable color so prefixes are plain ASCII and assertions are stable.
	prev := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = prev }()

	out := captureStderr(t, func() { FailMsg("error code %d: %s", 42, "bad input") })
	assert.Contains(t, out, "error code 42: bad input")
}
