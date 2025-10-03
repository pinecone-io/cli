package help

import (
	"strings"

	"github.com/MakeNowJust/heredoc"
)

const pad = "  "

// Examples normalizes a Cobra Example block
// - Accepts a multi-line string with possible indentation
// - De-indents with heredoc.Doc, trims leading/trailing whitespace, preserves interior blank lines
// - Left-indents each line and adds $ for command lines
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
			out = append(out, pad+trimmed)
			continue
		}

		// Command line
		out = append(out, pad+"$ "+trimmed)
	}

	return strings.Join(out, "\n")
}

func Short(shortDesc string) string {
	return strings.TrimSpace(heredoc.Doc(shortDesc))
}

func ShortF(shortDesc string, args ...any) string {
	return strings.TrimSpace(heredoc.Docf(shortDesc, args...))
}

func Long(longDesc string) string {
	return normalize(heredoc.Doc(longDesc))
}

func LongF(longDesc string, args ...any) string {
	return normalize(heredoc.Docf(longDesc, args...))
}

func normalize(s string) string {
	s = strings.TrimSpace(s)
	// Normalize CRLF -> LF
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return s
}
