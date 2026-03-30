package style

import (
	"os"

	"github.com/fatih/color"
	"golang.org/x/term"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
)

func colorEnabled() bool {
	return config.Color.Get() && term.IsTerminal(int(os.Stdout.Fd())) && term.IsTerminal(int(os.Stderr.Fd()))
}

func applyColor(s string, c *color.Color) string {
	color.NoColor = !colorEnabled()
	colored := c.SprintFunc()
	return colored(s)
}

func applyStyle(s string, c color.Attribute) string {
	color.NoColor = !colorEnabled()
	colored := color.New(c).SprintFunc()
	return colored(s)
}

func Faint(s string) string {
	return applyStyle(s, color.Faint)
}

func CodeWithPrompt(s string) string {
	return (applyStyle("$ ", color.Faint) + applyColor(s, color.New(color.FgMagenta, color.Bold)))
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
