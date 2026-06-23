// Package pluginhint emits a one-line marker that prompts Claude Code users to
// install the official Pinecone plugin.
//
// When the CLI detects it is running inside Claude Code (via the CLAUDECODE or
// CLAUDE_CODE_CHILD_SESSION environment variables), it writes a self-closing
// <claude-code-hint /> tag to stderr on its own line. Claude Code reads the
// marker, strips it from the command output before it reaches the model (so it
// never appears in the conversation or counts toward tokens), and shows the
// user a one-time prompt to install the plugin.
//
// The protocol is documented at https://code.claude.com/docs/en/plugin-hints.
package pluginhint

import (
	"io"
	"os"
	"sync"
)

// hintLine is the marker Claude Code looks for. It must occupy its own line and
// reference a plugin in an official Anthropic-controlled marketplace, otherwise
// Claude Code drops it.
//
//	v     — protocol version (1 is the only supported value)
//	type  — hint kind (plugin is the only supported value)
//	value — plugin identifier in name@marketplace form
const hintLine = `<claude-code-hint v="1" type="plugin" value="pinecone@claude-plugins-official" />`

// once guards against emitting the hint more than once per process. Claude Code
// deduplicates by plugin, so a single emission per invocation is all that's
// needed even though Emit may be reachable from multiple code paths.
var once sync.Once

// inClaudeCode reports whether the CLI appears to be running inside a Claude
// Code session.
//
//   - CLAUDECODE is set on every Claude Code version, giving the widest reach.
//   - CLAUDE_CODE_CHILD_SESSION is set (v2.1.172+) only in subprocesses Claude
//     Code itself spawns, such as Bash tool calls.
//
// Gating on either keeps the marker out of normal human-run invocations.
func inClaudeCode() bool {
	return os.Getenv("CLAUDECODE") != "" || os.Getenv("CLAUDE_CODE_CHILD_SESSION") != ""
}

// Emit writes the plugin hint to stderr when running inside Claude Code. It is
// safe to call from any number of code paths and on every invocation; the hint
// is written at most once per process and is a no-op outside Claude Code.
//
// stderr keeps the tag out of shell pipelines (e.g. `pc index list | jq`),
// though Claude Code scans both streams.
func Emit() {
	emitTo(os.Stderr)
}

// emitTo is the testable core of Emit.
func emitTo(w io.Writer) {
	if !inClaudeCode() {
		return
	}
	once.Do(func() {
		// Written on its own line; the trailing newline is required so the tag
		// occupies a line by itself.
		_, _ = io.WriteString(w, hintLine+"\n")
	})
}
