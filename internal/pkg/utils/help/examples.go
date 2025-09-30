package help

import (
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

const pad = "  "

// Examples normalizes a Cobra example block.
// - Accepts a multi-line string with possible indentation
// - De-indents with heredoc.Doc, trims leading/trailing whitespace, preserves interior blank lines
// - Left-indents each line and applies styles.CodeWithPrompt for command lines
func Examples(examples string) string {
	block := strings.TrimSpace(heredoc.Doc(examples))
	if block == "" {
		return ""
	}

	lines := strings.Split(block, "\n")
	out := make([]string, 0, len(lines))

	for _, line := range lines {
		line := strings.TrimRight(line, " \t\r")
		if strings.TrimSpace(line) == "" {
			out = append(out, "")
			continue
		}

		trimmed := strings.TrimLeft(line, " \t")

		// Comment line
		if strings.HasPrefix(trimmed, "#") {
			out = append(out, pad+style.Faint(trimmed))
			continue
		}

		// Command line
		out = append(out, pad+style.CodeWithPrompt(trimmed))
	}

	return strings.Join(out, "\n")
}
