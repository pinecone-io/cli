package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintFetchVectorsTable(resp *pinecone.FetchVectorsResponse) {
	writer := NewTabWriter()

	// Header Block
	if resp.Namespace != "" {
		pcio.Fprintf(writer, "Namespace: %s\n", resp.Namespace)
	}
	if resp.Usage != nil {
		pcio.Fprintf(writer, "Usage: %d (read units)\n", resp.Usage.ReadUnits)
	}

	// Table Header
	columns := []string{"ID", "DIMENSION", "VALUES", "SPARSE VALUES", "METADATA"}
	pcio.Fprintln(writer, strings.Join(columns, "\t"))

	// Rows
	for id, vector := range resp.Vectors {
		dim := 0
		if vector.Values != nil {
			dim = len(*vector.Values)
		}
		sparseDim := 0
		if vector.SparseValues != nil {
			sparseDim = len(vector.SparseValues.Values)
		}
		metadata := ""
		if vector.Metadata != nil {
			metadata = text.InlineJSON(vector.Metadata)
		}
		preview := previewSlice(vector.Values, 3)
		row := []string{id, pcio.Sprintf("%d", dim), preview, pcio.Sprintf("%d", sparseDim), metadata}
		pcio.Fprintln(writer, strings.Join(row, "\t"))
	}

	writer.Flush()
}

func previewSlice(values *[]float32, limit int) string {
	if values == nil || len(*values) == 0 {
		return "<none>"
	}
	vals := *values
	if len(vals) > limit {
		vals = vals[:limit]
	}
	return text.InlineJSON(vals) + "..."
}
