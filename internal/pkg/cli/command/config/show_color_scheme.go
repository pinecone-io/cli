package config

import (
	"fmt"

	conf "github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewShowColorSchemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-color-scheme",
		Short: "Display the Pinecone CLI color scheme for development reference",
		Long: `Display all available colors in the Pinecone CLI color scheme.
This command is useful for developers to see available colors and choose appropriate ones for their components.

Use 'pc config set-color-scheme' to change the color scheme.`,
		Example: help.Examples([]string{
			"pc config show-color-scheme",
		}),
		Run: func(cmd *cobra.Command, args []string) {
			showSimpleColorScheme()
		},
	}

	return cmd
}

// showSimpleColorScheme displays colors in a simple text format
func showSimpleColorScheme() {
	colorsEnabled := conf.Color.Get()

	fmt.Println("🎨 Pinecone CLI Color Scheme")
	fmt.Println("============================")
	fmt.Printf("Colors Enabled: %t\n", colorsEnabled)

	// Show which color scheme is being used
	currentScheme := conf.ColorScheme.Get()
	fmt.Printf("Color Scheme: %s\n", currentScheme)
	fmt.Println()

	if colorsEnabled {
		// Primary colors
		fmt.Println("Primary Colors:")
		fmt.Printf("  Primary Blue: %s\n", style.PrimaryStyle().Render("This is primary blue text"))
		fmt.Printf("  Success Green: %s\n", style.SuccessStyle().Render("This is success green text"))
		fmt.Printf("  Warning Yellow: %s\n", style.WarningStyle().Render("This is warning yellow text"))
		fmt.Printf("  Error Red: %s\n", style.ErrorStyle().Render("This is error red text"))
		fmt.Printf("  Info Blue: %s\n", style.InfoStyle().Render("This is info blue text"))
		fmt.Println()

		// Text colors
		fmt.Println("Text Colors:")
		fmt.Printf("  Primary Text: %s\n", style.PrimaryTextStyle().Render("This is primary text"))
		fmt.Printf("  Secondary Text: %s\n", style.SecondaryTextStyle().Render("This is secondary text"))
		fmt.Printf("  Muted Text: %s\n", style.MutedTextStyle().Render("This is muted text"))
		fmt.Println()

		// Background colors
		fmt.Println("Background Colors:")
		fmt.Printf("  Background: %s\n", style.BackgroundStyle().Render("This is background color"))
		fmt.Printf("  Surface: %s\n", style.SurfaceStyle().Render("This is surface color"))
		fmt.Println()

		// Border colors
		fmt.Println("Border Colors:")
		fmt.Printf("  Border: %s\n", style.BorderStyle().Render("This is border color"))
		fmt.Printf("  Border Muted: %s\n", style.BorderMutedStyle().Render("This is border muted color"))
		fmt.Println()

		// Usage examples with actual CLI function calls
		fmt.Println("Status Messages Examples:")
		fmt.Printf("  %s\n", style.SuccessMsg("Operation completed successfully"))
		fmt.Printf("  %s\n", style.FailMsg("Operation failed"))
		fmt.Printf("  %s\n", style.WarnMsg("This is a warning message"))
		fmt.Printf("  %s\n", style.InfoMsg("This is an info message"))
		fmt.Println()

		// Typography examples
		fmt.Println("Typography Examples:")
		fmt.Printf("  %s\n", style.Emphasis("This text is emphasized"))
		fmt.Printf("  %s\n", style.HeavyEmphasis("This text is heavily emphasized"))
		fmt.Printf("  %s\n", style.Heading("This is a heading"))
		fmt.Printf("  %s\n", style.Underline("This text is underlined"))
		fmt.Printf("  %s\n", style.Hint("This is a hint message"))
		fmt.Printf("  This is code/command: %s\n", style.Code("pc login"))
		fmt.Printf("  This is URL: %s\n", style.URL("https://pinecone.io"))
	} else {
		fmt.Println("Colors are disabled. Enable colors to see the color scheme.")
	}
}
