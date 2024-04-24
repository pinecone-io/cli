package presenters

import (
	"os"
	"text/tabwriter"
)

func NewTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 12, 1, 4, ' ', 0)
}
