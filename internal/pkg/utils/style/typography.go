package style

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
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

func Underline(s string) string {
	return applyStyle(s, color.Underline)
}

func Hint(s string) string {
	return applyStyle("Hint: ", color.Faint) + s
}

func CodeHint(templateString string, codeString string) string {
	return applyStyle("Hint: ", color.Faint) + pcio.Sprintf(templateString, Code(codeString))
}

func SuccessMsg(s string) string {
	return applyStyle("[SUCCESS] ", color.FgGreen) + s
}

func FailMsg(s string, a ...interface{}) string {
	return applyStyle("[ERROR] ", color.FgRed) + fmt.Sprintf(s, a...)
}

func Code(s string) string {
	formatted := applyColor(s, color.New(color.FgMagenta, color.Bold))
	if color.NoColor {
		// Add backticks for code formatting if color is disabled
		return "`" + formatted + "`"
	}
	return formatted
}

func URL(s string) string {
	return applyStyle(applyStyle(s, color.FgBlue), color.Italic)
}
