package presenters

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func NewTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
}

// PrintEmptyState prints a consistent placeholder for nil presenter inputs.
// It always returns true so callers can use it in guard clauses:
// if resp == nil && PrintEmptyState(writer, "vectors") { return }.
func PrintEmptyState(writer *tabwriter.Writer, resource string) bool {
	if resource == "" {
		resource = "data"
	}
	fmt.Fprintf(writer, "No %s available.\n", resource)
	writer.Flush()
	return true
}
