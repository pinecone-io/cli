package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/index"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

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
	var rows []presenters.Row
	for _, idx := range indexes {
		values := ExtractValuesForIndexAttributesGroups(idx, nonEmptyGroups)
		rows = append(rows, presenters.Row(values))
	}

	// Use the table utility
	presenters.PrintTable(presenters.TableOptions{
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
	var rows []presenters.Row
	for i, group := range nonEmptyGroups {
		// Get the columns with full names for this specific group
		groupColumns := getColumnsWithNamesForIndexAttributesGroup(group)
		groupValues := getValuesForIndexAttributesGroup(idx, group)

		// Add spacing before each group (except the first)
		if i > 0 {
			rows = append(rows, presenters.Row{"", ""})
		}

		// Add rows for this group using full names
		for j, col := range groupColumns {
			if j < len(groupValues) {
				rows = append(rows, presenters.Row{col.FullTitle, groupValues[j]})
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
	// Add spacing after the last row
	fmt.Println()
}

// PrintIndexCreateConfigTable creates and renders a table for index creation configuration
func PrintIndexCreateConfigTable(config *index.CreateOptions) {
	fmt.Println(style.Heading("Index Configuration"))
	fmt.Println()

	// Build rows for the table using the same order as the table view
	var rows []presenters.Row

	// Essential information
	rows = append(rows, presenters.Row{"Name", config.Name})
	rows = append(rows, presenters.Row{"Specification", config.GetSpecString()})

	// Vector type (for serverless)
	if config.VectorType != "" {
		rows = append(rows, presenters.Row{"Vector Type", config.VectorType})
	} else {
		rows = append(rows, presenters.Row{"Vector Type", "dense"}) // Default
	}

	rows = append(rows, presenters.Row{"Metric", config.Metric})

	if config.Dimension > 0 {
		rows = append(rows, presenters.Row{"Dimension", fmt.Sprintf("%d", config.Dimension)})
	}

	// Add spacing
	rows = append(rows, presenters.Row{"", ""})

	// Spec-specific information
	spec := config.GetSpecString()
	switch spec {
	case "serverless":
		rows = append(rows, presenters.Row{"Cloud Provider", config.Cloud})
		rows = append(rows, presenters.Row{"Region", config.Region})
	case "pod":
		rows = append(rows, presenters.Row{"Environment", config.Environment})
		rows = append(rows, presenters.Row{"Pod Type", config.PodType})
		rows = append(rows, presenters.Row{"Replicas", fmt.Sprintf("%d", config.Replicas)})
		rows = append(rows, presenters.Row{"Shard Count", fmt.Sprintf("%d", config.Shards)})
	}

	// Add spacing
	rows = append(rows, presenters.Row{"", ""})

	// Other information
	if config.DeletionProtection != "" {
		rows = append(rows, presenters.Row{"Deletion Protection", config.DeletionProtection})
	}

	if len(config.Tags) > 0 {
		var tagStrings []string
		for key, value := range config.Tags {
			tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", key, value))
		}
		rows = append(rows, presenters.Row{"Tags", strings.Join(tagStrings, ", ")})
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
	// Add spacing after the last row
	fmt.Println()
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
