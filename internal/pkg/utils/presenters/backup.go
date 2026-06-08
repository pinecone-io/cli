package presenters

import (
	"fmt"
	"strings"
	"time"

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
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Backup ID\t%s\n", backup.BackupId)
	fmt.Fprintf(writer, "Name\t%s\n", DisplayOrNone(backup.Name))
	fmt.Fprintf(writer, "Description\t%s\n", DisplayOrNone(backup.Description))
	fmt.Fprintf(writer, "Status\t%s\n", colorizeBackupStatus(backup.Status))
	fmt.Fprintf(writer, "Source Index\t%s\n", backup.SourceIndexName)
	fmt.Fprintf(writer, "Source Index ID\t%s\n", backup.SourceIndexId)
	fmt.Fprintf(writer, "Cloud\t%s\n", backup.Cloud)
	fmt.Fprintf(writer, "Region\t%s\n", backup.Region)
	fmt.Fprintf(writer, "Dimension\t%s\n", DisplayOrNone(backup.Dimension))
	fmt.Fprintf(writer, "Record Count\t%s\n", DisplayOrNone(backup.RecordCount))
	fmt.Fprintf(writer, "Namespace Count\t%s\n", DisplayOrNone(backup.NamespaceCount))
	fmt.Fprintf(writer, "Size (bytes)\t%s\n", DisplayOrNone(backup.SizeBytes))
	schema := "<none>"
	if backup.Schema != nil {
		schema = text.InlineJSON(backup.Schema)
	}
	fmt.Fprintf(writer, "Schema\t%s\n", schema)
	fmt.Fprintf(writer, "Created At\t%s\n", DisplayOrNone(backup.CreatedAt))
	fmt.Fprintf(writer, "Tags\t%s\n", formatTagsInline(backup.Tags))

	writer.Flush()
}

func PrintBackupList(list *pinecone.BackupList) {
	if list == nil || len(list.Data) == 0 {
		w := NewTabWriter()
		PrintEmptyState(w, "backups")
		return
	}

	cols := []tableColumn{
		{header: "BACKUP ID"},
		{header: "NAME"},
		{header: "INDEX"},
		{header: "STATUS", colorizer: colorizeBackupStatus},
		{header: "CLOUD/REGION"},
		{header: "RECORDS"},
		{header: "NAMESPACES"},
		{header: "SIZE (B)"},
		{header: "CREATED"},
	}
	rows := make([][]string, len(list.Data))
	for i, b := range list.Data {
		rows[i] = []string{
			b.BackupId,
			DisplayOrNone(b.Name),
			b.SourceIndexName,
			b.Status,
			fmt.Sprintf("%s/%s", b.Cloud, b.Region),
			DisplayOrNone(b.RecordCount),
			DisplayOrNone(b.NamespaceCount),
			DisplayOrNone(b.SizeBytes),
			DisplayOrNone(b.CreatedAt),
		}
	}
	printColorizedTable(cols, rows)

	if list.Pagination != nil && list.Pagination.Next != "" {
		fmt.Printf("\nNext Pagination Token: %s\n", list.Pagination.Next)
	}
}

func PrintRestoreJob(job *pinecone.RestoreJob) {
	writer := NewTabWriter()
	if job == nil {
		PrintEmptyState(writer, "restore job details")
		return
	}

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Restore Job ID\t%s\n", job.RestoreJobId)
	fmt.Fprintf(writer, "Backup ID\t%s\n", job.BackupId)
	fmt.Fprintf(writer, "Target Index\t%s\n", job.TargetIndexName)
	fmt.Fprintf(writer, "Status\t%s\n", colorizeRestoreJobStatus(job.Status))
	fmt.Fprintf(writer, "Percent Complete\t%s\n", DisplayOrNone(job.PercentComplete))
	fmt.Fprintf(writer, "Created At\t%s\n", formatTime(job.CreatedAt))
	fmt.Fprintf(writer, "Completed At\t%s\n", formatTimePtr(job.CompletedAt))

	writer.Flush()
}

func PrintRestoreJobList(list *pinecone.RestoreJobList) {
	if list == nil || len(list.Data) == 0 {
		w := NewTabWriter()
		PrintEmptyState(w, "restore jobs")
		return
	}

	cols := []tableColumn{
		{header: "RESTORE JOB ID"},
		{header: "BACKUP ID"},
		{header: "TARGET INDEX"},
		{header: "STATUS", colorizer: colorizeRestoreJobStatus},
		{header: "PERCENT"},
		{header: "CREATED"},
		{header: "COMPLETED"},
	}
	rows := make([][]string, len(list.Data))
	for i, job := range list.Data {
		rows[i] = []string{
			job.RestoreJobId,
			job.BackupId,
			job.TargetIndexName,
			job.Status,
			DisplayOrNone(job.PercentComplete),
			formatTime(job.CreatedAt),
			formatTimePtr(job.CompletedAt),
		}
	}
	printColorizedTable(cols, rows)

	if list.Pagination != nil && list.Pagination.Next != "" {
		fmt.Printf("\nNext Pagination Token: %s\n", list.Pagination.Next)
	}
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
