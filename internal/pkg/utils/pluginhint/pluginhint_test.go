package pluginhint

import (
	"bytes"
	"strings"
	"sync"
	"testing"
)

// resetOnce restores the package-level once guard so each test starts fresh.
func resetOnce() {
	once = sync.Once{}
}

func TestEmitTo_WritesHintWhenInClaudeCode(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
	}{
		{name: "CLAUDECODE set", env: map[string]string{"CLAUDECODE": "1"}},
		{name: "CLAUDE_CODE_CHILD_SESSION set", env: map[string]string{"CLAUDE_CODE_CHILD_SESSION": "1"}},
		{name: "both set", env: map[string]string{"CLAUDECODE": "1", "CLAUDE_CODE_CHILD_SESSION": "1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetOnce()
			t.Setenv("CLAUDECODE", "")
			t.Setenv("CLAUDE_CODE_CHILD_SESSION", "")
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			var buf bytes.Buffer
			emitTo(&buf)

			got := buf.String()
			if got != hintLine+"\n" {
				t.Errorf("emitTo() = %q, want %q", got, hintLine+"\n")
			}
		})
	}
}

func TestEmitTo_NoopOutsideClaudeCode(t *testing.T) {
	resetOnce()
	t.Setenv("CLAUDECODE", "")
	t.Setenv("CLAUDE_CODE_CHILD_SESSION", "")

	var buf bytes.Buffer
	emitTo(&buf)

	if buf.Len() != 0 {
		t.Errorf("emitTo() wrote %q outside Claude Code, want nothing", buf.String())
	}
}

func TestEmitTo_OnlyEmitsOnce(t *testing.T) {
	resetOnce()
	t.Setenv("CLAUDECODE", "1")

	var buf bytes.Buffer
	emitTo(&buf)
	emitTo(&buf)
	emitTo(&buf)

	if got := strings.Count(buf.String(), "claude-code-hint"); got != 1 {
		t.Errorf("emitTo() emitted hint %d times, want 1", got)
	}
}

func TestHintLine_MatchesProtocol(t *testing.T) {
	// The tag must be self-closing with the three required attributes, target
	// the official marketplace, and contain no newline (it must occupy a single
	// line; emitTo appends the terminating newline).
	if strings.Contains(hintLine, "\n") {
		t.Errorf("hintLine must not contain a newline: %q", hintLine)
	}
	for _, want := range []string{
		`v="1"`,
		`type="plugin"`,
		`value="pinecone@claude-plugins-official"`,
		"<claude-code-hint",
		"/>",
	} {
		if !strings.Contains(hintLine, want) {
			t.Errorf("hintLine missing %q; got %q", want, hintLine)
		}
	}
}
