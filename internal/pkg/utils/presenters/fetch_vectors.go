package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

type FetchVectorsResults struct {
	Vectors    map[string]*pinecone.Vector `json:"vectors,omitempty"`
	Namespace  string                      `json:"namespace"`
	Usage      *pinecone.Usage             `json:"usage,omitempty"`
	Pagination *pinecone.Pagination        `json:"pagination,omitempty"`
}

func NewFetchVectorsResultsFromFetch(resp *pinecone.FetchVectorsResponse) *FetchVectorsResults {
	if resp == nil {
		return &FetchVectorsResults{}
	}
	return &FetchVectorsResults{
		Vectors:   resp.Vectors,
		Namespace: resp.Namespace,
		Usage:     resp.Usage,
	}
}

func NewFetchVectorsResultsFromFetchByMetadata(resp *pinecone.FetchVectorsByMetadataResponse) *FetchVectorsResults {
	if resp == nil {
		return &FetchVectorsResults{}
	}
	return &FetchVectorsResults{
		Vectors:    resp.Vectors,
		Namespace:  resp.Namespace,
		Usage:      resp.Usage,
		Pagination: resp.Pagination,
	}
}

func PrintFetchVectorsTable(results *FetchVectorsResults) {
	writer := NewTabWriter()

	// Header Block
	if results.Namespace != "" {
		pcio.Fprintf(writer, "Namespace: %s\n", results.Namespace)
	}
	if results.Usage != nil {
		pcio.Fprintf(writer, "Usage: %d (read units)\n", results.Usage.ReadUnits)
	}

	// Table Header
	columns := []string{"ID", "DIMENSION", "VALUES", "SPARSE VALUES", "METADATA"}
	pcio.Fprintln(writer, strings.Join(columns, "\t"))

	// Rows
	for id, vector := range results.Vectors {
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
