package style

import (
	"github.com/fatih/color"
)

func Emphasis(s string) string {
	cyan := color.New(color.FgCyan).SprintFunc()
	return cyan(s)
}

func Code(s string) string {
	magenta := color.New(color.FgMagenta).SprintFunc()
	return magenta(s)
}

func StatusGreen(s string) string {
	green := color.New(color.FgGreen).SprintFunc()
	return green(s)
}

func StatusYellow(s string) string {
	yellow := color.New(color.FgYellow).SprintFunc()
	return yellow(s)
}

func StatusRed(s string) string {
	red := color.New(color.FgRed).SprintFunc()
	return red(s)
}
