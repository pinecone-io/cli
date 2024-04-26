package style

import (
	"github.com/fatih/color"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
)

func applyStyle(s string, c color.Attribute) string {
	color.NoColor = !config.Color.Get()
	colored := color.New(c).SprintFunc()
	return colored(s)
}

func Emphasis(s string) string {
	return applyStyle(s, color.FgCyan)
}

func Code(s string) string {
	formatted := applyStyle(s, color.FgMagenta)
	if color.NoColor {
		// Add backticks for code formatting if color is disabled
		return "`" + formatted + "`"
	}
	return formatted
}

func StatusGreen(s string) string {
	return applyStyle(s, color.FgGreen)
}

func StatusYellow(s string) string {
	return applyStyle(s, color.FgYellow)
}

func StatusRed(s string) string {
	return applyStyle(s, color.FgRed)
}
