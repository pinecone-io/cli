package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
)

// ColorScheme defines the centralized color palette for the Pinecone CLI
// Based on Pinecone's official website CSS variables for consistent branding
type ColorScheme struct {
	// Primary brand colors
	PrimaryBlue   string // Pinecone blue - main brand color
	SuccessGreen  string // Success states
	WarningYellow string // Warning states
	ErrorRed      string // Error states
	InfoBlue      string // Info states

	// Text colors
	PrimaryText   string // Main text color
	SecondaryText string // Secondary/muted text
	MutedText     string // Very muted text

	// Background colors
	Background string // Main background
	Surface    string // Surface/card backgrounds

	// Border colors
	Border      string // Default borders
	BorderMuted string // Muted borders
}

// Available color schemes
var AvailableColorSchemes = map[string]ColorScheme{
	"pc-default-dark":  DarkColorScheme(),
	"pc-default-light": LightColorScheme(),
}

// LightColorScheme returns colors optimized for light terminal backgrounds
// Uses Pinecone's official colors for text/backgrounds, but vibrant colors for status messages
func LightColorScheme() ColorScheme {
	return ColorScheme{
		// Primary brand colors (using vibrant colors that work well in both themes)
		PrimaryBlue:   "#002bff", // --primary-main (Pinecone brand)
		SuccessGreen:  "#28a745", // More vibrant green for better visibility
		WarningYellow: "#ffc107", // More vibrant amber for better visibility
		ErrorRed:      "#dc3545", // More vibrant red for better visibility
		InfoBlue:      "#17a2b8", // More vibrant info blue for better visibility

		// Text colors (from Pinecone's light theme)
		PrimaryText:   "#1c1917", // --text-primary
		SecondaryText: "#57534e", // --text-secondary
		MutedText:     "#a8a29e", // --text-tertiary

		// Background colors (from Pinecone's light theme)
		Background: "#fbfbfc", // --background
		Surface:    "#f2f3f6", // --surface

		// Border colors (from Pinecone's light theme)
		Border:      "#e7e5e4", // --border
		BorderMuted: "#d8dddf", // --divider
	}
}

// DarkColorScheme returns colors optimized for dark terminal backgrounds
// Uses Pinecone's official colors for text/backgrounds, but more vibrant colors for status messages
func DarkColorScheme() ColorScheme {
	return ColorScheme{
		// Primary brand colors (optimized for dark terminals)
		PrimaryBlue:   "#1e86ee", // --primary-main
		SuccessGreen:  "#28a745", // More vibrant green for dark terminals
		WarningYellow: "#ffc107", // More vibrant amber for dark terminals
		ErrorRed:      "#dc3545", // More vibrant red for dark terminals
		InfoBlue:      "#17a2b8", // More vibrant info blue for dark terminals

		// Text colors (from Pinecone's dark theme)
		PrimaryText:   "#fff",    // --text-primary
		SecondaryText: "#a3a3a3", // --text-secondary
		MutedText:     "#525252", // --text-tertiary

		// Background colors (from Pinecone's dark theme)
		Background: "#171717", // --background
		Surface:    "#252525", // --surface

		// Border colors (from Pinecone's dark theme)
		Border:      "#404040", // --border
		BorderMuted: "#2a2a2a", // --divider
	}
}

// DefaultColorScheme returns the configured color scheme
func DefaultColorScheme() ColorScheme {
	schemeName := config.ColorScheme.Get()
	if scheme, exists := AvailableColorSchemes[schemeName]; exists {
		return scheme
	}
	// Fallback to dark theme if configured scheme doesn't exist
	return DarkColorScheme()
}

// GetColorScheme returns the current color scheme
// This can be extended in the future to support themes
func GetColorScheme() ColorScheme {
	return DefaultColorScheme()
}

// LipglossColorScheme provides lipgloss-compatible color styles
type LipglossColorScheme struct {
	PrimaryBlue   lipgloss.Color
	SuccessGreen  lipgloss.Color
	WarningYellow lipgloss.Color
	ErrorRed      lipgloss.Color
	InfoBlue      lipgloss.Color
	PrimaryText   lipgloss.Color
	SecondaryText lipgloss.Color
	MutedText     lipgloss.Color
	Background    lipgloss.Color
	Surface       lipgloss.Color
	Border        lipgloss.Color
	BorderMuted   lipgloss.Color
}

// GetLipglossColorScheme returns lipgloss-compatible colors
func GetLipglossColorScheme() LipglossColorScheme {
	scheme := GetColorScheme()
	return LipglossColorScheme{
		PrimaryBlue:   lipgloss.Color(scheme.PrimaryBlue),
		SuccessGreen:  lipgloss.Color(scheme.SuccessGreen),
		WarningYellow: lipgloss.Color(scheme.WarningYellow),
		ErrorRed:      lipgloss.Color(scheme.ErrorRed),
		InfoBlue:      lipgloss.Color(scheme.InfoBlue),
		PrimaryText:   lipgloss.Color(scheme.PrimaryText),
		SecondaryText: lipgloss.Color(scheme.SecondaryText),
		MutedText:     lipgloss.Color(scheme.MutedText),
		Background:    lipgloss.Color(scheme.Background),
		Surface:       lipgloss.Color(scheme.Surface),
		Border:        lipgloss.Color(scheme.Border),
		BorderMuted:   lipgloss.Color(scheme.BorderMuted),
	}
}

// GetAvailableColorSchemeNames returns a list of available color scheme names
func GetAvailableColorSchemeNames() []string {
	names := make([]string, 0, len(AvailableColorSchemes))
	for name := range AvailableColorSchemes {
		names = append(names, name)
	}
	return names
}
