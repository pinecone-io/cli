package text

import (
	"strings"
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
