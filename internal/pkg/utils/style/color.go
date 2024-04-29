package style

import (
	"github.com/fatih/color"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
)

func applyColor(s string, c *color.Color) string {
	color.NoColor = !config.Color.Get()
	colored := c.SprintFunc()
	return colored(s)
}

func applyStyle(s string, c color.Attribute) string {
	color.NoColor = !config.Color.Get()
	colored := color.New(c).SprintFunc()
	return colored(s)
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
