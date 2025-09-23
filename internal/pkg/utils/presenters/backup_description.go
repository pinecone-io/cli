package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

func PrintDescribeBackupTable(backup *pinecone.Backup) {
	writer := NewTabWriter()
	log.Debug().Str("id", backup.BackupId).Msg("Printing backup description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	backupName := "unnamed"
	if backup.Name != nil {
		backupName = *backup.Name
	}
	pcio.Fprintf(writer, "Name\t%s\n", backupName)
	pcio.Fprintf(writer, "ID\t%s\n", backup.BackupId)
	pcio.Fprintf(writer, "Index Name\t%s\n", backup.SourceIndexName)
	pcio.Fprintf(writer, "Status\t%s\n", ColorizeBackupStatus(backup.Status))

	if backup.CreatedAt != nil {
		pcio.Fprintf(writer, "Created At\t%s\n", *backup.CreatedAt)
	} else {
		pcio.Fprintf(writer, "Created At\t<none>\n")
	}

	if backup.SizeBytes != nil {
		pcio.Fprintf(writer, "Size\t%d\n", *backup.SizeBytes)
	} else {
		pcio.Fprintf(writer, "Size\t<none>\n")
	}

	writer.Flush()
}

func ColorizeBackupStatus(status string) string {
	switch status {
	case "Ready":
		return style.SuccessStyle().Render(status)
	case "InProgress":
		return style.WarningStyle().Render(status)
	case "Failed":
		return style.ErrorStyle().Render(status)
	}

	return status
}
