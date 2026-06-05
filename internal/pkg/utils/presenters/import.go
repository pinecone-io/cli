package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

// PrintStartImportTable prints a summary of a newly started import operation.
func PrintStartImportTable(resp *pinecone.StartImportResponse) {
	writer := NewTabWriter()
	if resp == nil {
		PrintEmptyState(writer, "import details")
		return
	}

	columns := []string{"ATTRIBUTE", "VALUE"}
	fmt.Fprint(writer, strings.Join(columns, "\t")+"\n")
	fmt.Fprintf(writer, "Import ID\t%s\n", resp.Id)

	writer.Flush()
}

// PrintImportTable prints the full details of a single import operation.
func PrintImportTable(imp *pinecone.Import) {
	writer := NewTabWriter()
	if imp == nil {
		PrintEmptyState(writer, "import details")
		return
	}

	columns := []string{"ATTRIBUTE", "VALUE"}
	fmt.Fprint(writer, strings.Join(columns, "\t")+"\n")

	fmt.Fprintf(writer, "Import ID\t%s\n", imp.Id)
	fmt.Fprintf(writer, "Status\t%s\n", colorizeImportStatus(imp.Status))
	fmt.Fprintf(writer, "URI\t%s\n", imp.Uri)
	fmt.Fprintf(writer, "Percent Complete\t%.1f%%\n", imp.PercentComplete)
	fmt.Fprintf(writer, "Records Imported\t%d\n", imp.RecordsImported)
	fmt.Fprintf(writer, "Created At\t%s\n", formatTimePtr(imp.CreatedAt))
	fmt.Fprintf(writer, "Finished At\t%s\n", formatTimePtr(imp.FinishedAt))
	fmt.Fprintf(writer, "Error\t%s\n", DisplayOrNone(imp.Error))

	writer.Flush()
}

// PrintImportList prints a table of import operations.
func PrintImportList(list *pinecone.ListImportsResponse) {
	writer := NewTabWriter()
	if list == nil || len(list.Imports) == 0 {
		PrintEmptyState(writer, "imports")
		return
	}

	columns := []string{"IMPORT ID", "STATUS", "URI", "PERCENT", "RECORDS", "CREATED", "FINISHED"}
	fmt.Fprint(writer, strings.Join(columns, "\t")+"\n")

	for _, imp := range list.Imports {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%.1f%%\t%d\t%s\t%s\n",
			imp.Id,
			colorizeImportStatus(imp.Status),
			imp.Uri,
			imp.PercentComplete,
			imp.RecordsImported,
			formatTimePtr(imp.CreatedAt),
			formatTimePtr(imp.FinishedAt),
		)
	}

	if list.NextPaginationToken != nil && *list.NextPaginationToken != "" {
		fmt.Fprintf(writer, "\nNext Pagination Token: %s\n", *list.NextPaginationToken)
	}

	writer.Flush()
}

func colorizeImportStatus(status pinecone.ImportStatus) string {
	switch status {
	case pinecone.Completed:
		return style.StatusGreen(string(status))
	case pinecone.InProgress, pinecone.Pending:
		return style.StatusYellow(string(status))
	case pinecone.Failed, pinecone.Cancelled:
		return style.StatusRed(string(status))
	default:
		return string(status)
	}
}
