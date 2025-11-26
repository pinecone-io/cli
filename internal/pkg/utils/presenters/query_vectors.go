package presenters

import (
	"sort"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintQueryVectorsTable(resp *pinecone.QueryVectorsResponse) {
	writer := NewTabWriter()
	if resp == nil {
		PrintEmptyState(writer, "query results")
		return
	}

	// Header Block
	if resp.Namespace != "" {
		pcio.Fprintf(writer, "Namespace: %s\n", resp.Namespace)
	}
	if resp.Usage != nil {
		pcio.Fprintf(writer, "Usage: %d (read units)\n", resp.Usage.ReadUnits)
	}

	// Detect which columns to show
	hasDense := false
	hasSparse := false
	hasMetadata := false
	for _, m := range resp.Matches {
		if m == nil || m.Vector == nil {
			continue
		}
		if m.Vector.Values != nil && len(*m.Vector.Values) > 0 {
			hasDense = true
		}
		if m.Vector.SparseValues != nil &&
			(len(m.Vector.SparseValues.Indices) > 0 || len(m.Vector.SparseValues.Values) > 0) {
			hasSparse = true
		}
		if m.Vector.Metadata != nil {
			hasMetadata = true
		}
	}

	// Table Header
	cols := []string{"ID", "SCORE"}
	if hasDense {
		cols = append(cols, "VALUES")
	}
	if hasSparse {
		cols = append(cols, "SPARSE INDICES", "SPARSE VALUES")
	}
	if hasMetadata {
		cols = append(cols, "METADATA")
	}
	pcio.Fprintln(writer, strings.Join(cols, "\t"))

	// Rows
	for _, match := range resp.Matches {
		if match == nil || match.Vector == nil {
			continue
		}
		row := []string{match.Vector.Id, pcio.Sprintf("%f", match.Score)}

		if hasDense {
			values := "<none>"
			if match.Vector.Values != nil {
				values = previewSliceFloat32(match.Vector.Values, 3)
			}
			row = append(row, values)
		}
		if hasSparse {
			iPreview, vPreview := "<none>", "<none>"
			if match.Vector.SparseValues != nil {
				iPreview = previewSliceUint32(match.Vector.SparseValues.Indices, 3)
				vPreview = previewSliceFloat32(&match.Vector.SparseValues.Values, 3)
			}
			row = append(row, iPreview, vPreview)
		}
		if hasMetadata {
			metadata := "<none>"
			if match.Vector.Metadata != nil {
				m := match.Vector.Metadata.AsMap()
				if len(m) > 0 {
					keys := make([]string, 0, len(m))
					for k := range m {
						keys = append(keys, k)
					}
					sort.Strings(keys)
					show := keys
					if len(show) > 3 {
						show = show[:3]
					}
					limited := make(map[string]any, len(show))
					for _, k := range show {
						limited[k] = m[k]
					}

					s := text.InlineJSON(limited) // compact one-line JSON
					if len(keys) > 3 {
						// put ellipsis inside the braces: {"a":1,"b":2,"c":3, ...}
						s = strings.TrimRight(s, "}") + ", ...}"
					}
					metadata = s
				}
			}
			row = append(row, metadata)
		}

		pcio.Fprintln(writer, strings.Join(row, "\t"))
	}

	writer.Flush()
}

func previewSliceUint32(values []uint32, limit int) string {
	if len(values) == 0 {
		return "<none>"
	}
	vals := values
	truncated := false
	if len(vals) > limit {
		vals = vals[:limit]
		truncated = true
	}

	text := text.InlineJSON(vals)
	if truncated && strings.HasSuffix(text, "]") {
		text = text[:len(text)-1] + ", ...]"
	}

	return text
}
