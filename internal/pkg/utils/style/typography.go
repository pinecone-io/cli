package style

import (
	"github.com/fatih/color"
)

func Emphasis(s string) string {
	return applyStyle(s, color.FgCyan)
}

func HeavyEmphasis(s string) string {
	return applyColor(s, color.New(color.FgCyan, color.Bold))
}

func Heading(s string) string {
	return applyStyle(s, color.Bold)
}

func Code(s string) string {
	formatted := applyStyle(s, color.FgMagenta)
	if color.NoColor {
		// Add backticks for code formatting if color is disabled
		return "`" + formatted + "`"
	}
	return formatted
}

func URL(s string) string {
	return applyStyle(applyStyle(s, color.FgBlue), color.Italic)
}
