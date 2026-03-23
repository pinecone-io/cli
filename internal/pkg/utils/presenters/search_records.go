package presenters

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintSearchRecordsTable(resp *pinecone.SearchRecordsResponse) {
	writer := NewTabWriter()
	if resp == nil {
		PrintEmptyState(writer, "search results")
		return
	}

	if resp.Usage.ReadUnits > 0 {
		fmt.Fprintf(writer, "Usage: %d (read units)\n", resp.Usage.ReadUnits)
	}
	if resp.Usage.EmbedTotalTokens != nil {
		fmt.Fprintf(writer, "Embed tokens: %d\n", *resp.Usage.EmbedTotalTokens)
	}
	if resp.Usage.RerankUnits != nil {
		fmt.Fprintf(writer, "Rerank units: %d\n", *resp.Usage.RerankUnits)
	}

	fmt.Fprintln(writer, "ID\tSCORE\tFIELDS")

	for _, hit := range resp.Result.Hits {
		fields := previewFields(hit.Fields, 3)
		fmt.Fprintf(writer, "%s\t%f\t%s\n", hit.Id, hit.Score, fields)
	}

	writer.Flush()
}

func previewFields(fields map[string]any, limit int) string {
	if len(fields) == 0 {
		return "<none>"
	}

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	show := keys
	truncated := false
	if len(show) > limit {
		show = show[:limit]
		truncated = true
	}

	limited := make(map[string]any, len(show))
	for _, k := range show {
		limited[k] = fields[k]
	}

	out := text.InlineJSON(limited)
	if truncated && strings.HasSuffix(out, "}") {
		out = strings.TrimSuffix(out, "}") + ", ...}"
	}

	return out
}
