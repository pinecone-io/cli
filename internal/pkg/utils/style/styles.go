package style

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
)

const MessageBoxFixedWidth = 14

// Predefined styles for common use cases
var (
	// Status styles
	SuccessStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.SuccessGreen)
	}

	WarningStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.WarningYellow)
	}

	ErrorStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.ErrorRed)
	}

	InfoStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.InfoBlue)
	}

	PrimaryStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.PrimaryBlue)
	}

	// Text styles
	PrimaryTextStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.PrimaryText)
	}

	SecondaryTextStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.SecondaryText)
	}

	MutedTextStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.MutedText)
	}

	// Background styles
	BackgroundStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Background(colors.Background).Foreground(colors.PrimaryText)
	}

	SurfaceStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Background(colors.Surface).Foreground(colors.PrimaryText)
	}

	// Border styles
	BorderStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.Border)
	}

	BorderMutedStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.BorderMuted)
	}

	// Typography styles
	EmphasisStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.PrimaryBlue)
	}

	HeavyEmphasisStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.PrimaryBlue).Bold(true)
	}

	HeadingStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.PrimaryText).Bold(true)
	}

	UnderlineStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.PrimaryText).Underline(true)
	}

	HintStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.SecondaryText)
	}

	CodeStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.InfoBlue).Bold(true)
	}

	URLStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().Foreground(colors.InfoBlue).Italic(true)
	}

	// Message box styles with icon|label layout
	SuccessBoxStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().
			Background(colors.SuccessGreen).
			Foreground(lipgloss.Color("#FFFFFF")). // Always white text for good contrast
			Bold(true).
			Padding(0, 1).
			Width(MessageBoxFixedWidth) // Fixed width for consistent alignment
	}

	WarningBoxStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().
			Background(colors.WarningYellow).
			Foreground(lipgloss.Color("#000000")). // Always black text for good contrast on yellow
			Bold(true).
			Padding(0, 1).
			Width(MessageBoxFixedWidth) // Fixed width for consistent alignment
	}

	ErrorBoxStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().
			Background(colors.ErrorRed).
			Foreground(lipgloss.Color("#FFFFFF")). // Always white text for good contrast
			Bold(true).
			Padding(0, 1).
			Width(MessageBoxFixedWidth) // Fixed width for consistent alignment
	}

	InfoBoxStyle = func() lipgloss.Style {
		colors := GetLipglossColorScheme()
		return lipgloss.NewStyle().
			Background(colors.InfoBlue).
			Foreground(lipgloss.Color("#FFFFFF")). // Always white text for good contrast
			Bold(true).
			Padding(0, 1).
			Width(MessageBoxFixedWidth) // Fixed width for consistent alignment
	}
)

// GetBrandedTableStyles returns table styles using the centralized color scheme
func GetBrandedTableStyles() (table.Styles, bool) {
	colors := GetLipglossColorScheme()
	colorsEnabled := config.Color.Get()

	s := table.DefaultStyles()

	if colorsEnabled {
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colors.PrimaryBlue).
			Foreground(colors.PrimaryBlue).
			BorderBottom(true).
			Bold(true)
		s.Cell = s.Cell.Padding(0, 1)
		// Ensure selected row style doesn't interfere
		s.Selected = s.Selected.
			Foreground(colors.PrimaryText).
			Background(colors.Background).
			Bold(false)
	} else {
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			Bold(true)
		s.Cell = s.Cell.Padding(0, 1)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("")).
			Background(lipgloss.Color("")).
			Bold(false)
	}

	return s, colorsEnabled
}

// GetBrandedTableNoSelectionStyles returns table styles for read-only tables without row selection
func GetBrandedTableNoSelectionStyles() (table.Styles, bool) {
	colors := GetLipglossColorScheme()
	colorsEnabled := config.Color.Get()

	s := table.DefaultStyles()

	if colorsEnabled {
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colors.PrimaryBlue).
			Foreground(colors.PrimaryBlue).
			BorderBottom(true).
			Bold(true)
		s.Cell = s.Cell.Padding(0, 1)
		// Empty selected style since cell style is already applied to each cell
		// and we don't want any additional styling for selected rows
		s.Selected = lipgloss.NewStyle()
	} else {
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			Bold(true)
		s.Cell = s.Cell.Padding(0, 1)
		// Empty selected style since cell style is already applied to each cell
		// and we don't want any additional styling for selected rows
		s.Selected = lipgloss.NewStyle()
	}

	return s, colorsEnabled
}

// GetInteractiveStyles returns styles for interactive components using the centralized color scheme
func GetInteractiveStyles() (lipgloss.Style, lipgloss.Style, lipgloss.Style, bool) {
	colors := GetLipglossColorScheme()
	colorsEnabled := config.Color.Get()

	var questionStyle, promptStyle, keyStyle lipgloss.Style

	if colorsEnabled {
		questionStyle = HeadingStyle()
		promptStyle = SecondaryTextStyle().MarginBottom(1)
		keyStyle = lipgloss.NewStyle().
			Foreground(colors.InfoBlue).
			Bold(true)
	} else {
		questionStyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1)

		promptStyle = lipgloss.NewStyle().
			MarginBottom(1)

		keyStyle = lipgloss.NewStyle().
			Bold(true)
	}

	return questionStyle, promptStyle, keyStyle, colorsEnabled
}
