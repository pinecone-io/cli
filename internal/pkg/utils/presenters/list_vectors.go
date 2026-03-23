package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintListVectorsTable(resp *pinecone.ListVectorsResponse) {
	writer := NewTabWriter()
	if resp == nil {
		PrintEmptyState(writer, "vector IDs")
		return
	}

	// Header block
	if resp.Namespace != "" {
		fmt.Fprintf(writer, "Namespace: %s\n", resp.Namespace)
	}
	if resp.Usage != nil {
		fmt.Fprintf(writer, "Usage: %d (read units)\n", resp.Usage.ReadUnits)
	}

	// Table header
	columns := []string{"ID"}
	fmt.Fprintln(writer, strings.Join(columns, "\t"))

	// Rows
	for _, vectorId := range resp.VectorIds {
		id := ""
		if vectorId != nil {
			id = *vectorId
		}
		fmt.Fprintln(writer, id)
	}

	// Pagination footer
	if resp.NextPaginationToken != nil {
		fmt.Fprintf(writer, "Next pagination token: %s\n", *resp.NextPaginationToken)
	}

	writer.Flush()
}
