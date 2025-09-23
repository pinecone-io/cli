package presenters

import (
	"fmt"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

// BackupDisplayData represents the unified display structure for backup information
type BackupDisplayData struct {
	// Essential information
	Name            string
	ID              string
	SourceIndexName string
	Status          string

	// Metadata information
	CreatedAt   string
	SizeBytes   string
	Description string

	// Source information
	SourceIndexId string
	Cloud         string
	Region        string

	// Technical information
	Dimension      string
	Metric         string
	RecordCount    string
	NamespaceCount string

	// Other information
	Tags string
}

// ConvertBackupToDisplayData converts a pinecone.Backup to BackupDisplayData
func ConvertBackupToDisplayData(backup *pinecone.Backup) *BackupDisplayData {
	data := &BackupDisplayData{}

	// Essential information
	if backup.Name != nil {
		data.Name = *backup.Name
	} else {
		data.Name = "unnamed"
	}
	data.ID = backup.BackupId
	data.SourceIndexName = backup.SourceIndexName
	data.Status = backup.Status

	// Metadata information
	if backup.CreatedAt != nil {
		data.CreatedAt = *backup.CreatedAt
	}
	if backup.SizeBytes != nil {
		data.SizeBytes = fmt.Sprintf("%d", *backup.SizeBytes)
	}
	if backup.Description != nil {
		data.Description = *backup.Description
	}

	// Source information
	data.SourceIndexId = backup.SourceIndexId
	data.Cloud = backup.Cloud
	data.Region = backup.Region

	// Technical information
	if backup.Dimension != nil {
		data.Dimension = fmt.Sprintf("%d", *backup.Dimension)
	}
	if backup.Metric != nil {
		data.Metric = string(*backup.Metric)
	}
	if backup.RecordCount != nil {
		data.RecordCount = fmt.Sprintf("%d", *backup.RecordCount)
	}
	if backup.NamespaceCount != nil {
		data.NamespaceCount = fmt.Sprintf("%d", *backup.NamespaceCount)
	}

	// Other information
	if backup.Tags != nil {
		data.Tags = presenters.FormatTags(backup.Tags)
	}

	return data
}

// PrintBackupDisplayTable creates and renders a table for backup display data
func PrintBackupDisplayTable(data *BackupDisplayData) {
	// Build rows for the table
	var rows []presenters.Row

	// Essential information
	rows = append(rows, presenters.Row{"Name", data.Name})
	rows = append(rows, presenters.Row{"ID", data.ID})
	rows = append(rows, presenters.Row{"Source Index", data.SourceIndexName})
	rows = append(rows, presenters.Row{"Status", presenters.ColorizeStatus(data.Status)})

	// Add spacing
	rows = append(rows, presenters.Row{"", ""})

	// Metadata information
	if data.CreatedAt != "" {
		rows = append(rows, presenters.Row{"Created At", data.CreatedAt})
	}
	if data.SizeBytes != "" {
		rows = append(rows, presenters.Row{"Size", data.SizeBytes})
	}
	if data.Description != "" {
		rows = append(rows, presenters.Row{"Description", data.Description})
	}

	// Add spacing
	rows = append(rows, presenters.Row{"", ""})

	// Source information
	rows = append(rows, presenters.Row{"Source Index ID", data.SourceIndexId})
	rows = append(rows, presenters.Row{"Cloud Provider", data.Cloud})
	rows = append(rows, presenters.Row{"Region", data.Region})

	// Add spacing
	rows = append(rows, presenters.Row{"", ""})

	// Technical information
	if data.Dimension != "" {
		rows = append(rows, presenters.Row{"Dimension", data.Dimension})
	}
	if data.Metric != "" {
		rows = append(rows, presenters.Row{"Metric", data.Metric})
	}
	if data.RecordCount != "" {
		rows = append(rows, presenters.Row{"Record Count", data.RecordCount})
	}
	if data.NamespaceCount != "" {
		rows = append(rows, presenters.Row{"Namespace Count", data.NamespaceCount})
	}

	// Other information
	if data.Tags != "" {
		rows = append(rows, presenters.Row{"", ""})
		rows = append(rows, presenters.Row{"Tags", data.Tags})
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

// PrintDescribeBackupTable creates and renders a table for backup description with right-aligned first column and secondary text styling
func PrintDescribeBackupTable(backup *pinecone.Backup) {
	// Print title
	fmt.Println(style.Heading("Backup Configuration"))
	fmt.Println()

	// Convert to display data and print
	data := ConvertBackupToDisplayData(backup)
	PrintBackupDisplayTable(data)
}

// PrintBackupTable creates and renders a table for backup list with proper column formatting
func PrintBackupTable(backups []*pinecone.Backup) {
	if len(backups) == 0 {
		fmt.Println("No backups found.")
		return
	}

	// Define table columns
	columns := []presenters.Column{
		{Title: "NAME", Width: 20},
		{Title: "ID", Width: 40},
		{Title: "SOURCE INDEX", Width: 20},
		{Title: "STATUS", Width: 12},
		{Title: "CREATED", Width: 16},
		{Title: "SIZE", Width: 8},
	}

	// Convert backups to table rows
	rows := make([]presenters.Row, len(backups))
	for i, backup := range backups {
		backupName := "unnamed"
		if backup.Name != nil {
			backupName = *backup.Name
		}

		created := "-"
		if backup.CreatedAt != nil {
			// Try to parse the timestamp and format it nicely
			if parsedTime, err := time.Parse(time.RFC3339, *backup.CreatedAt); err == nil {
				created = parsedTime.Format("Jan 02 15:04")
			} else {
				// If parsing fails, use the raw value
				created = *backup.CreatedAt
			}
		}

		size := "-"
		if backup.SizeBytes != nil {
			size = presenters.FormatSize(*backup.SizeBytes)
		}

		rows[i] = presenters.Row{
			backupName,
			backup.BackupId,
			backup.SourceIndexName,
			backup.Status, // Use unstyled status for table
			created,
			size,
		}
	}

	// Print the table
	presenters.PrintTable(presenters.TableOptions{
		Columns: columns,
		Rows:    rows,
	})

	fmt.Println()

	// Add a note about full details
	hint := fmt.Sprintf("Use %s to see backup details", style.Code("pc backup describe --id <backup-id>"))
	fmt.Println(style.Hint(hint))
}
