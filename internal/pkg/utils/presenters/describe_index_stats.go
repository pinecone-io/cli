package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintDescribeIndexStatsTable(resp *pinecone.DescribeIndexStatsResponse) {
	writer := NewTabWriter()

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Dimension\t%d\n", resp.Dimension)
	pcio.Fprintf(writer, "Index Fullness\t%f\n", resp.IndexFullness)
	pcio.Fprintf(writer, "Total Vector Count\t%d\n", resp.TotalVectorCount)

	formatted := text.IndentJSON(resp.Namespaces)
	formatted = strings.ReplaceAll(formatted, "\n", "\n\t") // indent lines under the namespace value
	pcio.Fprintf(writer, "Namespaces\t%s\n", formatted)

	writer.Flush()
}
