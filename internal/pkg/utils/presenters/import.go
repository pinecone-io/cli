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
	fmt.Fprintf(writer, "Status\t%s\n", colorizeImportStatus(string(imp.Status)))
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
	if list == nil || len(list.Imports) == 0 {
		w := NewTabWriter()
		PrintEmptyState(w, "imports")
		return
	}

	cols := []tableColumn{
		{header: "IMPORT ID"},
		{header: "STATUS", colorizer: colorizeImportStatus},
		{header: "URI"},
		{header: "PERCENT"},
		{header: "RECORDS"},
		{header: "CREATED"},
		{header: "FINISHED"},
	}
	rows := make([][]string, len(list.Imports))
	for i, imp := range list.Imports {
		rows[i] = []string{
			imp.Id,
			string(imp.Status),
			imp.Uri,
			fmt.Sprintf("%.1f%%", imp.PercentComplete),
			fmt.Sprintf("%d", imp.RecordsImported),
			formatTimePtr(imp.CreatedAt),
			formatTimePtr(imp.FinishedAt),
		}
	}
	printColorizedTable(cols, rows)

	if list.NextPaginationToken != nil && *list.NextPaginationToken != "" {
		fmt.Printf("\nNext Pagination Token: %s\n", *list.NextPaginationToken)
	}
}

func colorizeImportStatus(status string) string {
	switch pinecone.ImportStatus(status) {
	case pinecone.Completed:
		return style.StatusGreen(status)
	case pinecone.InProgress, pinecone.Pending:
		return style.StatusYellow(status)
	case pinecone.Failed, pinecone.Cancelled:
		return style.StatusRed(status)
	default:
		return status
	}
}
