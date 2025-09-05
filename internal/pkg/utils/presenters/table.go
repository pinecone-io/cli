// Package presenters provides table rendering functions for data display.
//
// NOTE: This package uses fmt functions directly (not pcio) because:
// - Data output should NOT be suppressed by the -q flag
// - Informational commands (list, describe) need to display data even in quiet mode
// - Only user-facing messages (progress, confirmations) should respect quiet mode
package presenters

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

// Column represents a table column with title and width
type Column struct {
	Title string
	Width int
}

// Row represents a table row as a slice of strings
type Row []string

// TableOptions contains configuration options for creating a table
type TableOptions struct {
	Columns []Column
	Rows    []Row
}

// PrintTable creates and renders a bubbles table with the given options
func PrintTable(options TableOptions) {
	// Convert abstract types to bubbles table types
	bubblesColumns := make([]table.Column, len(options.Columns))
	for i, col := range options.Columns {
		bubblesColumns[i] = table.Column{
			Title: col.Title,
			Width: col.Width,
		}
	}

	bubblesRows := make([]table.Row, len(options.Rows))
	for i, row := range options.Rows {
		bubblesRows[i] = table.Row(row)
	}

	// Create and configure the table
	t := table.New(
		table.WithColumns(bubblesColumns),
		table.WithRows(bubblesRows),
		table.WithFocused(false), // Always disable focus to prevent row selection
		table.WithHeight(len(options.Rows)),
	)

	// Use centralized color scheme for table styling (no selection version)
	s, _ := style.GetBrandedTableNoSelectionStyles()
	t.SetStyles(s)

	// Always ensure no row is selected/highlighted
	// This must be done after setting styles
	t.SetCursor(-1)

	// Render the table directly
	fmt.Println(t.View())
}

// PrintTableWithTitle creates and renders a bubbles table with a title
func PrintTableWithTitle(title string, options TableOptions) {
	fmt.Println()
	fmt.Printf("%s\n\n", style.Heading(title))
	PrintTable(options)
	fmt.Println()
}
