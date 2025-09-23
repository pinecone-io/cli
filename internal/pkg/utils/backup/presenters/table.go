package presenters

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

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
		{Title: "CREATED", Width: 25},
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
			created = presenters.FormatDate(*backup.CreatedAt)
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
	hint := fmt.Sprintf("Use %s to see backup details", style.Code("pc backup describe <backup-id>"))
	fmt.Println(style.Hint(hint))
}
