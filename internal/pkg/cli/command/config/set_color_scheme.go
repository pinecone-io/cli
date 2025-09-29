package config

import (
	"strings"

	conf "github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewSetColorSchemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-color-scheme",
		Short: "Configure the color scheme for the Pinecone CLI",
		Long: `Set the color scheme used by the Pinecone CLI.

Available color schemes:
  pc-default-dark  - Dark theme optimized for dark terminal backgrounds
  pc-default-light - Light theme optimized for light terminal backgrounds

The color scheme affects all colored output in the CLI, including tables, messages, and the color scheme display.`,
		Example: help.Examples([]string{
			"pc config set-color-scheme pc-default-dark",
			"pc config set-color-scheme pc-default-light",
		}),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				msg.FailMsg("Please provide a color scheme name")
				msg.InfoMsg("Available color schemes: %s", strings.Join(getAvailableColorSchemes(), ", "))
				exit.ErrorMsg("No color scheme provided")
			}

			schemeName := args[0]

			// Validate the color scheme
			if !isValidColorScheme(schemeName) {
				msg.FailMsg("Invalid color scheme: %s", schemeName)
				msg.InfoMsg("Available color schemes: %s", strings.Join(getAvailableColorSchemes(), ", "))
				exit.ErrorMsg("Invalid color scheme")
			}

			conf.ColorScheme.Set(schemeName)
			msg.SuccessMsg("Color scheme updated to %s", style.Emphasis(schemeName))
		},
	}

	return cmd
}

// getAvailableColorSchemes returns a list of available color scheme names
func getAvailableColorSchemes() []string {
	schemes := make([]string, 0, len(style.AvailableColorSchemes))
	for name := range style.AvailableColorSchemes {
		schemes = append(schemes, name)
	}
	return schemes
}

// isValidColorScheme checks if the given scheme name is valid
func isValidColorScheme(schemeName string) bool {
	_, exists := style.AvailableColorSchemes[schemeName]
	return exists
}
