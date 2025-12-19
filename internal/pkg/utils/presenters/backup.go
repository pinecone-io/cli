package presenters

import (
	"strings"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintBackupTable(backup *pinecone.Backup) {
	writer := NewTabWriter()
	if backup == nil {
		PrintEmptyState(writer, "backup details")
		return
	}

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Backup ID\t%s\n", backup.BackupId)
	pcio.Fprintf(writer, "Name\t%s\n", DisplayOrNone(backup.Name))
	pcio.Fprintf(writer, "Status\t%s\n", colorizeBackupStatus(backup.Status))
	pcio.Fprintf(writer, "Source Index\t%s\n", backup.SourceIndexName)
	pcio.Fprintf(writer, "Cloud\t%s\n", backup.Cloud)
	pcio.Fprintf(writer, "Region\t%s\n", backup.Region)
	pcio.Fprintf(writer, "Record Count\t%s\n", DisplayOrNone(backup.RecordCount))
	pcio.Fprintf(writer, "Namespace Count\t%s\n", DisplayOrNone(backup.NamespaceCount))
	pcio.Fprintf(writer, "Size (bytes)\t%s\n", DisplayOrNone(backup.SizeBytes))
	pcio.Fprintf(writer, "Metric\t%s\n", DisplayOrNone(backup.Metric))
	schema := "<none>"
	if backup.Schema != nil {
		schema = text.InlineJSON(backup.Schema)
	}
	pcio.Fprintf(writer, "Schema\t%s\n", schema)
	pcio.Fprintf(writer, "Created At\t%s\n", DisplayOrNone(backup.CreatedAt))
	pcio.Fprintf(writer, "Tags\t%s\n", formatTagsInline(backup.Tags))

	writer.Flush()
}

func PrintBackupList(list *pinecone.BackupList) {
	writer := NewTabWriter()
	if list == nil || len(list.Data) == 0 {
		PrintEmptyState(writer, "backups")
		return
	}

	columns := []string{"BACKUP ID", "NAME", "INDEX", "STATUS", "CLOUD/REGION", "RECORDS", "NAMESPACES", "SIZE (B)", "CREATED"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, b := range list.Data {
		cloudRegion := pcio.Sprintf("%s/%s", b.Cloud, b.Region)
		created := DisplayOrNone(b.CreatedAt)
		pcio.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			b.BackupId,
			DisplayOrNone(b.Name),
			b.SourceIndexName,
			b.Status,
			cloudRegion,
			DisplayOrNone(b.RecordCount),
			DisplayOrNone(b.NamespaceCount),
			DisplayOrNone(b.SizeBytes),
			created,
		)
	}

	if list.Pagination != nil && list.Pagination.Next != "" {
		pcio.Fprintf(writer, "\nNext Pagination Token: %s\n", list.Pagination.Next)
	}

	writer.Flush()
}

func PrintRestoreJob(job *pinecone.RestoreJob) {
	writer := NewTabWriter()
	if job == nil {
		PrintEmptyState(writer, "restore job details")
		return
	}

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Restore Job ID\t%s\n", job.RestoreJobId)
	pcio.Fprintf(writer, "Backup ID\t%s\n", job.BackupId)
	pcio.Fprintf(writer, "Target Index\t%s\n", job.TargetIndexName)
	pcio.Fprintf(writer, "Status\t%s\n", colorizeRestoreJobStatus(job.Status))
	pcio.Fprintf(writer, "Percent Complete\t%s\n", DisplayOrNone(job.PercentComplete))
	pcio.Fprintf(writer, "Created At\t%s\n", formatTime(job.CreatedAt))
	pcio.Fprintf(writer, "Completed At\t%s\n", formatTimePtr(job.CompletedAt))

	writer.Flush()
}

func PrintRestoreJobList(list *pinecone.RestoreJobList) {
	writer := NewTabWriter()
	if list == nil || len(list.Data) == 0 {
		PrintEmptyState(writer, "restore jobs")
		return
	}

	columns := []string{"RESTORE JOB ID", "BACKUP ID", "TARGET INDEX", "STATUS", "PERCENT", "CREATED", "COMPLETED"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, job := range list.Data {
		pcio.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			job.RestoreJobId,
			job.BackupId,
			job.TargetIndexName,
			job.Status,
			DisplayOrNone(job.PercentComplete),
			formatTime(job.CreatedAt),
			formatTimePtr(job.CompletedAt),
		)
	}

	if list.Pagination != nil && list.Pagination.Next != "" {
		pcio.Fprintf(writer, "\nNext Pagination Token: %s\n", list.Pagination.Next)
	}

	writer.Flush()
}

func colorizeBackupStatus(status string) string {
	switch strings.ToLower(status) {
	case "ready":
		return style.StatusGreen(status)
	case "initializing":
		return style.StatusYellow(status)
	case "initializationfailed":
		return style.StatusRed(status)
	default:
		return status
	}
}

func colorizeRestoreJobStatus(status string) string {
	switch strings.ToLower(status) {
	case "completed":
		return style.StatusGreen(status)
	case "pending":
		return style.StatusYellow(status)
	case "failed", "cancelled":
		return style.StatusRed(status)
	default:
		return status
	}
}

func formatTagsInline(tags *pinecone.IndexTags) string {
	if tags == nil {
		return "<none>"
	}
	return text.InlineJSON(tags)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "<none>"
	}
	return t.UTC().Format(time.RFC3339)
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return "<none>"
	}
	return formatTime(*t)
}
