package presenters

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func NewTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
}

// tableColumn describes a single column in a printColorizedTable call.
type tableColumn struct {
	header    string
	colorizer func(string) string // nil = no colorization
}

// printColorizedTable renders a tab-aligned table to stdout. Rows contain plain
// string values; each column's optional colorizer is applied after tabwriter has
// computed column widths, so ANSI bytes never affect alignment.
// Columns are processed right-to-left so an earlier substitution never shifts
// the byte offsets of columns already replaced in the same line.
func printColorizedTable(cols []tableColumn, rows [][]string) {
	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = c.header
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 12, 1, 3, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	fmt.Println(lines[0])

	// Locate each column's byte offset in the formatted header line.
	offsets := make([]int, len(cols))
	cursor := 0
	for i, h := range headers {
		pos := strings.Index(lines[0][cursor:], h)
		if pos < 0 {
			continue
		}
		offsets[i] = cursor + pos
		cursor = offsets[i] + len(h)
	}

	for i, row := range rows {
		line := lines[i+1]
		for j := len(cols) - 1; j >= 0; j-- {
			if cols[j].colorizer == nil {
				continue
			}
			plain := row[j]
			colored := cols[j].colorizer(plain)
			if colored == plain {
				continue
			}
			start := offsets[j]
			if start+len(plain) <= len(line) {
				line = line[:start] + colored + line[start+len(plain):]
			}
		}
		fmt.Println(line)
	}
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
