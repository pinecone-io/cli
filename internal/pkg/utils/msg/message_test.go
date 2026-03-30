package msg

import (
	"bytes"
	"io"
	"os"
	"strings"
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

// captureStdout redirects os.Stdout to a pipe for the duration of f,
// returning everything written to stdout as a string.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	prev := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = prev }()

	f()
	w.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("captureStdout: reading pipe: %v", err)
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

func TestFailJSON_StripANSIFromJSONOutput(t *testing.T) {
	// Simulate the --json pipeline case: stderr is a TTY (colors enabled) but
	// the JSON value on stdout must still be clean plain text.
	prev := color.NoColor
	color.NoColor = false // force colors on, as if stderr were a TTY
	defer func() { color.NoColor = prev }()

	var stdout string
	captureStderr(t, func() {
		stdout = captureStdout(t, func() {
			FailJSON(true, "key %s not found", "\x1b[36mmy-key\x1b[0m")
		})
	})

	assert.NotContains(t, stdout, "\x1b[", "JSON output must not contain ANSI escape codes")
	assert.Contains(t, stdout, "my-key", "JSON output must contain the plain text value")
}

func TestFailJSON_WithJSONFlag(t *testing.T) {
	prev := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = prev }()

	var stdout, stderr string
	stderr = captureStderr(t, func() {
		stdout = captureStdout(t, func() {
			FailJSON(true, "failed to list indexes: %s", "timeout")
		})
	})

	// stdout should contain a JSON error object
	assert.Contains(t, stdout, `"error"`)
	assert.Contains(t, stdout, "failed to list indexes: timeout")
	assert.True(t, strings.HasPrefix(strings.TrimSpace(stdout), "{"), "stdout should be a JSON object")

	// stderr should still contain the human-readable message
	assert.Contains(t, stderr, "failed to list indexes: timeout")
}

func TestFailJSON_WithoutJSONFlag(t *testing.T) {
	prev := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = prev }()

	var stdout, stderr string
	stderr = captureStderr(t, func() {
		stdout = captureStdout(t, func() {
			FailJSON(false, "something broke")
		})
	})

	// stdout should be empty
	assert.Empty(t, stdout)

	// stderr should contain the human-readable message
	assert.Contains(t, stderr, "something broke")
}

func TestFailJSON_NoANSIOnNonTTY(t *testing.T) {
	// With color.NoColor = true (non-TTY), output should contain no ANSI escape sequences.
	prev := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = prev }()

	out := captureStderr(t, func() { FailMsg("plain error") })
	assert.NotContains(t, out, "\x1b[")
}
