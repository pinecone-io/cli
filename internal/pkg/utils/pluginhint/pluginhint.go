// Package pluginhint implements the Claude Code plugin-hint protocol: when the
// CLI runs inside Claude Code, it writes a marker to stderr that prompts the
// user to install the official Pinecone plugin. Claude Code strips the marker
// from the output before it reaches the model and shows a one-time install
// prompt. See https://code.claude.com/docs/en/plugin-hints.
package pluginhint

import (
	"io"
	"os"
	"sync"
)

// hintLine must occupy its own line and reference a plugin in an official
// Anthropic marketplace; Claude Code drops hints that don't.
const hintLine = `<claude-code-hint v="1" type="plugin" value="pinecone@claude-plugins-official" />`

var once sync.Once

// CLAUDECODE is set on every Claude Code version; CLAUDE_CODE_CHILD_SESSION
// (v2.1.172+) only in subprocesses it spawns, such as Bash tool calls.
func inClaudeCode() bool {
	return os.Getenv("CLAUDECODE") != "" || os.Getenv("CLAUDE_CODE_CHILD_SESSION") != ""
}

// Emit writes the plugin hint to stderr when running inside Claude Code. It is
// a no-op otherwise, and emits at most once per process, so it is safe to call
// on every invocation and from multiple code paths.
func Emit() {
	emitTo(os.Stderr)
}

func emitTo(w io.Writer) {
	if !inClaudeCode() {
		return
	}
	once.Do(func() {
		_, _ = io.WriteString(w, hintLine+"\n")
	})
}
