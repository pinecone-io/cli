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
