package help

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	cobra "github.com/spf13/cobra"
)

// Colors help examples with ANSI color codes only when color is enabled
// Allows stripping ANSI color codes when generating man pages
func renderExamples(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}

	lines := strings.Split(s, "\n")
	for i, ln := range lines {
		raw := strings.TrimLeft(ln, " ")
		left := ln[:len(ln)-len(raw)]

		// Comment line
		if strings.HasPrefix(raw, "# ") {
			lines[i] = left + style.Faint(raw)
			continue
		}

		// Command line
		if strings.HasPrefix(raw, "$ ") {
			// Color without duplicating the $ prefix
			cmd := strings.TrimPrefix(raw, "$ ")
			lines[i] = left + style.CodeWithPrompt(cmd)
		}
		// Other lines - leave as is
	}

	return strings.Join(lines, "\n")
}

func EnableColorizedHelp(root *cobra.Command) {
	cobra.AddTemplateFunc("pcExamples", renderExamples)
}
