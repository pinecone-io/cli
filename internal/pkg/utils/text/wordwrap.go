package text

import (
	"strings"
	"unicode"
)

func WordWrap(text string, maxWidth int) string {
	// First, remove all remaining newlines
	text = strings.ReplaceAll(text, "\n", " ")

	var wrappedLines []string
	words := strings.Fields(text)
	line := ""
	lineLength := 0

	for _, word := range words {
		wordLength := len(word)

		if lineLength+wordLength <= maxWidth {
			line += word + " "
			lineLength += wordLength + 1
		} else {
			wrappedLines = append(wrappedLines, line)
			line = word + " "
			lineLength = wordLength + 1
		}
	}

	if line != "" {
		wrappedLines = append(wrappedLines, line)
	}

	return strings.Join(wrappedLines, "\n")
}

func WordWrapPreserveFormatting(text string, maxWidth int) string {
	var wrappedText strings.Builder
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		wrappedLine := lineWrap(line, maxWidth)
		wrappedText.WriteString(wrappedLine + "\n")
	}

	return strings.TrimSuffix(wrappedText.String(), "\n")
}

func lineWrap(line string, width int) string {
	var wrappedLine strings.Builder
	words := strings.Fields(line)
	lineLength := 0
	var leadingWhitespace string

	// Capture leading whitespace
	for _, r := range line {
		if unicode.IsSpace(r) {
			leadingWhitespace += string(r)
		} else {
			break
		}
	}

	for _, word := range words {
		if lineLength+len(word)+1 > width {
			wrappedLine.WriteString("\n" + leadingWhitespace)
			lineLength = len(leadingWhitespace)
		}
		if lineLength > 0 && lineLength != len(leadingWhitespace) {
			wrappedLine.WriteString(" ")
			lineLength++
		}
		wrappedLine.WriteString(word)
		lineLength += len(word)
	}

	return wrappedLine.String()
}
