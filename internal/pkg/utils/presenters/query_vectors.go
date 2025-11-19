package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintQueryVectorsTable(resp *pinecone.QueryVectorsResponse) {
	writer := NewTabWriter()

	// Header Block
	if resp.Namespace != "" {
		pcio.Fprintf(writer, "Namespace: %s\n", resp.Namespace)
	}
	if resp.Usage != nil {
		pcio.Fprintf(writer, "Usage: %d (read units)\n", resp.Usage.ReadUnits)
	}

	// Table Header
	columns := []string{"ID", "SCORE", "VALUES", "SPARSE INDICES", "SPARSE VALUES", "METADATA"}
	pcio.Fprintln(writer, strings.Join(columns, "\t"))

	// Rows
	for _, match := range resp.Matches {
		if match == nil || match.Vector == nil {
			continue
		}

		valuesPreview := "<none>"
		if match.Vector.Values != nil {
			valuesPreview = previewSliceFloat32(match.Vector.Values, 3)
		}

		sparseIndicesPreview := "<none>"
		if match.Vector.SparseValues != nil {
			sparseIndicesPreview = previewSliceUint32(match.Vector.SparseValues.Indices, 3)
		}

		sparseValuesPreview := "<none>"
		if match.Vector.SparseValues != nil {
			sparseValuesPreview = previewSliceFloat32(&match.Vector.SparseValues.Values, 3)
		}

		metadataPreview := "<none>"
		if match.Vector.Metadata != nil {
			metadataPreview = text.InlineJSON(match.Vector.Metadata)
		}

		row := []string{match.Vector.Id, pcio.Sprintf("%f", match.Score), valuesPreview, sparseIndicesPreview, sparseValuesPreview, metadataPreview}
		pcio.Fprintln(writer, strings.Join(row, "\t"))
	}
}

func previewSliceUint32(values []uint32, limit int) string {
	if len(values) == 0 {
		return "<none>"
	}
	vals := values
	if len(vals) > limit {
		vals = vals[:limit]
	}
	return text.InlineJSON(vals) + "..."
}
