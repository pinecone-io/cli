package help

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	defaultWidth = 80
	minWidth     = 40
	maxWidth     = 90
)

func EnableColorizedHelp(root *cobra.Command) {
	cobra.AddTemplateFunc("pcExamples", renderExamples)
	cobra.AddTemplateFunc("pcBlock", renderLongBlock)
}

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

func renderLongBlock(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}

	// Select a width based on the terminal, if we can
	width := resolveWrapWidth()

	lines := strings.Split(s, "\n")
	var out []string
	var para []string
	inCodeFence := false

	flush := func() {
		if len(para) == 0 {
			return
		}
		block := strings.Join(para, "\n")
		out = append(out, renderParagraph(block, width))
		para = para[:0]
	}

	for _, ln := range lines {
		trim := strings.TrimSpace(ln)

		if strings.HasPrefix(trim, "```") {
			inCodeFence = !inCodeFence
			para = append(para, ln)
			continue
		}

		if !inCodeFence && trim == "" {
			flush()
			out = append(out, "")
			continue
		}

		para = append(para, ln)
	}
	flush()

	return strings.Join(out, "\n")
}

func renderParagraph(para string, width int) string {
	lines := strings.Split(para, "\n")

	// We don't wrap fenced blocks, indented code blocks, or headings
	if isFenced(lines) || hasLeadingIndent(lines) || isListBlock(lines) || isHeadingBlock(lines) || hasCommandLine(lines) {
		return strings.Join(lines, "\n")
	}

	return wrap(strings.Join(lines, " "), width)
}

func isFenced(lines []string) bool {
	if len(lines) == 0 {
		return false
	}
	head := strings.TrimSpace(lines[0])
	tail := strings.TrimSpace(lines[len(lines)-1])
	return strings.HasPrefix(head, "```") && strings.HasPrefix(tail, "```")
}

func hasLeadingIndent(lines []string) bool {
	for _, ln := range lines {
		if ln == "" {
			continue
		}
		if ln[0] == ' ' || ln[0] == '\t' {
			return true
		}
	}
	return false
}

var reOrderedList = regexp.MustCompile(`^\d+\.\s+`)

func isListLine(s string) bool {
	t := strings.TrimLeft(s, " ")
	return strings.HasPrefix(t, "- ") || strings.HasPrefix(t, "* ") || reOrderedList.MatchString(t)
}

func isListBlock(lines []string) bool {
	any := false
	for _, ln := range lines {
		if ln == "" {
			continue
		}
		if !isListLine(ln) {
			return false
		}
		any = true
	}
	return any
}

func isHeadingLine(s string) bool {
	t := strings.TrimLeft(s, " ")
	return strings.HasPrefix(t, "# ")
}

func isHeadingBlock(lines []string) bool {
	return len(lines) == 1 && isHeadingLine(lines[0])
}

func hasCommandLine(lines []string) bool {
	for _, ln := range lines {
		if strings.HasPrefix(strings.TrimLeft(ln, " "), "$ ") {
			return true
		}
	}
	return false
}

func wrap(s string, width int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if width <= minWidth {
		width = defaultWidth
	}
	words := strings.Fields(s)
	var b strings.Builder
	lineLen := 0
	for i, w := range words {
		wlen := utf8.RuneCountInString(w)
		separator := 0
		if i > 0 {
			separator = 1
		}
		if lineLen+wlen+separator > width {
			b.WriteByte('\n')
			b.WriteString(w)
			lineLen = wlen
		} else {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(w)
			lineLen += separator + wlen
		}
	}
	return b.String()
}

func resolveWrapWidth() int {
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		return clamp(w-2, minWidth, maxWidth)
	}
	if w, _, err := term.GetSize(int(os.Stderr.Fd())); err == nil && w > 0 {
		return clamp(w-2, minWidth, maxWidth)
	}
	if w := atoiEnv("COLUMNS"); w > 0 {
		return clamp(w-2, minWidth, maxWidth)
	}
	return defaultWidth
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func atoiEnv(key string) int {
	if s := os.Getenv(key); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			return n
		}
	}
	return 0
}
