package style

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

// Typography functions using predefined styles

func Emphasis(s string) string {
	return EmphasisStyle().Render(s)
}

func HeavyEmphasis(s string) string {
	return HeavyEmphasisStyle().Render(s)
}

func Heading(s string) string {
	return HeadingStyle().Render(s)
}

func Underline(s string) string {
	return UnderlineStyle().Render(s)
}

func Hint(s string) string {
	return HintStyle().Render("Hint: ") + s
}

func CodeHint(templateString string, codeString string) string {
	return HintStyle().Render("Hint: ") + pcio.Sprintf(templateString, Code(codeString))
}

func Code(s string) string {
	if color.NoColor {
		// Add backticks for code formatting if color is disabled
		return "`" + s + "`"
	}
	return CodeStyle().Render(s)
}

func URL(s string) string {
	return URLStyle().Render(s)
}

// Message functions using predefined box styles

func SuccessMsg(s string) string {
	if color.NoColor {
		return fmt.Sprintf("🟩 [SUCCESS] %s", s)
	}
	icon := "\r🟩"
	box := SuccessBoxStyle().Render(icon + " SUCCESS")
	return fmt.Sprintf("%s %s", box, s)
}

func WarnMsg(s string) string {
	if color.NoColor {
		return fmt.Sprintf("🟧 [WARNING] %s", s)
	}
	icon := "\r🟧"
	box := WarningBoxStyle().Render(icon + " WARNING")
	return fmt.Sprintf("%s %s", box, s)
}

func InfoMsg(s string) string {
	if color.NoColor {
		return fmt.Sprintf("🟦 [INFO] %s", s)
	}
	icon := "\r🟦"
	box := InfoBoxStyle().Render(icon + " INFO")
	return fmt.Sprintf("%s %s", box, s)
}

func FailMsg(s string, a ...any) string {
	message := fmt.Sprintf(s, a...)
	if color.NoColor {
		return fmt.Sprintf("🟥 [ERROR] %s", message)
	}
	icon := "\r🟥"
	box := ErrorBoxStyle().Render(icon + " ERROR")
	return fmt.Sprintf("%s %s", box, message)
}

// repeat creates a string by repeating a character n times
func repeat(char string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += char
	}
	return result
}

// WarnMsgMultiLine creates a multi-line warning message with proper alignment
func WarnMsgMultiLine(messages ...string) string {
	if len(messages) == 0 {
		return ""
	}

	if color.NoColor {
		// Simple text format for no-color mode
		result := "🟧 [WARNING] " + messages[0]
		for i := 1; i < len(messages); i++ {
			result += "\n           " + messages[i]
		}
		return result
	}

	// Create the first line with the warning label
	icon := "\r🟧"
	label := " WARNING"
	boxStyle := WarningBoxStyle()
	labelBox := boxStyle.Render(icon + label)

	// Build the result
	result := labelBox + " " + messages[0]
	for i := 1; i < len(messages); i++ {
		result += "\n" + repeat(" ", MessageBoxFixedWidth) + messages[i]
	}

	return result
}

// SuccessMsgMultiLine creates a multi-line success message with proper alignment
func SuccessMsgMultiLine(messages ...string) string {
	if len(messages) == 0 {
		return ""
	}

	if color.NoColor {
		// Simple text format for no-color mode
		result := "🟩 [SUCCESS] " + messages[0]
		for i := 1; i < len(messages); i++ {
			result += "\n           " + messages[i]
		}
		return result
	}

	// Create the first line with the success label
	icon := "\r🟩"
	label := " SUCCESS"
	boxStyle := SuccessBoxStyle()
	labelBox := boxStyle.Render(icon + label)

	// Build the result
	result := labelBox + " " + messages[0]
	for i := 1; i < len(messages); i++ {
		result += "\n" + repeat(" ", MessageBoxFixedWidth) + messages[i]
	}

	return result
}

// InfoMsgMultiLine creates a multi-line info message with proper alignment
func InfoMsgMultiLine(messages ...string) string {
	if len(messages) == 0 {
		return ""
	}

	if color.NoColor {
		// Simple text format for no-color mode
		result := "🟦 [INFO] " + messages[0]
		for i := 1; i < len(messages); i++ {
			result += "\n          " + messages[i]
		}
		return result
	}

	// Create the first line with the info label
	icon := "\r🟦"
	label := " INFO"
	boxStyle := InfoBoxStyle()
	labelBox := boxStyle.Render(icon + label)

	// Build the result
	result := labelBox + " " + messages[0]
	for i := 1; i < len(messages); i++ {
		result += "\n" + repeat(" ", MessageBoxFixedWidth) + messages[i]
	}

	return result
}

// FailMsgMultiLine creates a multi-line error message with proper alignment
func FailMsgMultiLine(messages ...string) string {
	if len(messages) == 0 {
		return ""
	}

	if color.NoColor {
		// Simple text format for no-color mode
		result := "🟥 [ERROR] " + messages[0]
		for i := 1; i < len(messages); i++ {
			result += "\n           " + messages[i]
		}
		return result
	}

	// Create the first line with the error label
	icon := "\r🟥"
	label := " ERROR"
	boxStyle := ErrorBoxStyle()
	labelBox := boxStyle.Render(icon + label)

	// Build the result
	result := labelBox + " " + messages[0]
	for i := 1; i < len(messages); i++ {
		result += "\n" + repeat(" ", MessageBoxFixedWidth) + messages[i]
	}

	return result
}

// Legacy functions using fatih/color (kept for backward compatibility)

func CodeWithPrompt(s string) string {
	if color.NoColor {
		return "$ " + s
	}
	colors := GetLipglossColorScheme()
	promptStyle := lipgloss.NewStyle().Foreground(colors.SecondaryText)
	commandStyle := lipgloss.NewStyle().Foreground(colors.InfoBlue).Bold(true)
	return promptStyle.Render("$ ") + commandStyle.Render(s)
}
