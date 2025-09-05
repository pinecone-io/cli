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
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
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

// PrintIndexTableWithIndexAttributesGroups creates and renders a table for index information with custom index attribute groups
func PrintIndexTableWithIndexAttributesGroups(indexes []*pinecone.Index, groups []IndexAttributesGroup) {
	// Filter out groups that have no meaningful data
	nonEmptyGroups := filterNonEmptyIndexAttributesGroups(indexes, groups)
	if len(nonEmptyGroups) == 0 {
		return
	}

	// Get columns for the non-empty groups
	columns := GetColumnsForIndexAttributesGroups(nonEmptyGroups)

	// Build table rows
	var rows []Row
	for _, idx := range indexes {
		values := ExtractValuesForIndexAttributesGroups(idx, nonEmptyGroups)
		rows = append(rows, Row(values))
	}

	// Use the table utility
	PrintTable(TableOptions{
		Columns: columns,
		Rows:    rows,
	})

	fmt.Println()

	// Add a note about full URLs if state info is shown
	hasStateGroup := false
	for _, group := range nonEmptyGroups {
		if group == IndexAttributesGroupState {
			hasStateGroup = true
			break
		}
	}
	if hasStateGroup && len(indexes) > 0 {
		hint := fmt.Sprintf("Use %s to see index details", style.Code("pc index describe <name>"))
		fmt.Println(style.Hint(hint))
	}
}

// PrintDescribeIndexTable creates and renders a table for index description with right-aligned first column and secondary text styling
func PrintDescribeIndexTable(idx *pinecone.Index) {
	log.Debug().Str("name", idx.Name).Msg("Printing index description")

	// Print title
	fmt.Println(style.Heading("Index Configuration"))
	fmt.Println()

	// Print all groups with their information
	PrintDescribeIndexTableWithIndexAttributesGroups(idx, AllIndexAttributesGroups())
}

// PrintDescribeIndexTableWithIndexAttributesGroups creates and renders a table for index description with specified index attribute groups
func PrintDescribeIndexTableWithIndexAttributesGroups(idx *pinecone.Index, groups []IndexAttributesGroup) {
	// Filter out groups that have no meaningful data for this specific index
	nonEmptyGroups := filterNonEmptyIndexAttributesGroupsForIndex(idx, groups)
	if len(nonEmptyGroups) == 0 {
		return
	}

	// Build rows for the table using the same order as the table view
	var rows []Row
	for i, group := range nonEmptyGroups {
		// Get the columns with full names for this specific group
		groupColumns := getColumnsWithNamesForIndexAttributesGroup(group)
		groupValues := getValuesForIndexAttributesGroup(idx, group)

		// Add spacing before each group (except the first)
		if i > 0 {
			rows = append(rows, Row{"", ""})
		}

		// Add rows for this group using full names
		for j, col := range groupColumns {
			if j < len(groupValues) {
				rows = append(rows, Row{col.FullTitle, groupValues[j]})
			}
		}
	}

	// Print each row with right-aligned first column and secondary text styling
	for _, row := range rows {
		if len(row) >= 2 {
			// Right align the first column content
			rightAlignedFirstCol := fmt.Sprintf("%20s", row[0])

			// Apply secondary text styling to the first column
			styledFirstCol := style.SecondaryTextStyle().Render(rightAlignedFirstCol)

			// Print the row
			rowText := fmt.Sprintf("%s  %s", styledFirstCol, row[1])
			fmt.Println(rowText)
		} else if len(row) == 1 && row[0] == "" {
			// Empty row for spacing
			fmt.Println()
		}
	}
}

// ColorizeState applies appropriate styling to index state
func ColorizeState(state pinecone.IndexStatusState) string {
	switch state {
	case pinecone.Ready:
		return style.SuccessStyle().Render(string(state))
	case pinecone.Initializing, pinecone.Terminating, pinecone.ScalingDown, pinecone.ScalingDownPodSize, pinecone.ScalingUp, pinecone.ScalingUpPodSize:
		return style.WarningStyle().Render(string(state))
	case pinecone.InitializationFailed:
		return style.ErrorStyle().Render(string(state))
	default:
		return string(state)
	}
}

// ColorizeDeletionProtection applies appropriate styling to deletion protection status
func ColorizeDeletionProtection(deletionProtection pinecone.DeletionProtection) string {
	if deletionProtection == pinecone.DeletionProtectionEnabled {
		return style.SuccessStyle().Render("enabled")
	}
	return style.ErrorStyle().Render("disabled")
}
