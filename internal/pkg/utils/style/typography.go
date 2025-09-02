package style

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

// Lipgloss styles for cli-alerts style messages
var (
	// Alert type boxes (solid colored backgrounds) - using standard CLI colors
	successBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#28a745")). // Standard green
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1)

	errorBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#dc3545")). // Standard red (softer)
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1)

	warningBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#ffc107")). // Standard amber/yellow
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)

	infoBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#17a2b8")). // Standard info blue
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1)
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
	if color.NoColor {
		return fmt.Sprintf("✔ [SUCCESS] %s", s)
	}
	icon := "✔"
	box := successBoxStyle.Render(icon + " SUCCESS")
	return fmt.Sprintf("%s %s", box, s)
}

func WarnMsg(s string) string {
	if color.NoColor {
		return fmt.Sprintf("⚠ [WARNING] %s", s)
	}
	icon := "⚠"
	box := warningBoxStyle.Render(icon + " WARNING")
	return fmt.Sprintf("%s %s", box, s)
}

func InfoMsg(s string) string {
	if color.NoColor {
		return fmt.Sprintf("ℹ [INFO] %s", s)
	}
	icon := "ℹ"
	box := infoBoxStyle.Render(icon + " INFO")
	return fmt.Sprintf("%s %s", box, s)
}

func FailMsg(s string, a ...any) string {
	message := fmt.Sprintf(s, a...)
	if color.NoColor {
		return fmt.Sprintf("✘ [ERROR] %s", message)
	}
	icon := "✘"
	box := errorBoxStyle.Render(icon + " ERROR")
	return fmt.Sprintf("%s %s", box, message)
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
