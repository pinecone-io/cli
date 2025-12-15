package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintDescribeIndexStatsTable(resp *pinecone.DescribeIndexStatsResponse) {
	writer := NewTabWriter()
	if resp == nil {
		PrintEmptyState(writer, "index stats")
		return
	}

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	dimension := uint32(0)
	if resp.Dimension != nil {
		dimension = *resp.Dimension
	}

	pcio.Fprintf(writer, "Dimension\t%d\n", dimension)
	pcio.Fprintf(writer, "Index Fullness\t%f\n", resp.IndexFullness)
	pcio.Fprintf(writer, "Total Vector Count\t%d\n", resp.TotalVectorCount)

	if len(resp.Namespaces) == 0 {
		pcio.Fprintf(writer, "Namespaces\t<none>\n")
	} else {
		pcio.Fprintf(writer, "Namespaces\n")
		pcio.Fprintf(writer, "\tNAME\tVECTOR COUNT\n")

		names := make([]string, 0, len(resp.Namespaces))
		for name := range resp.Namespaces {
			names = append(names, name)
		}
		for _, name := range names {
			pcio.Fprintf(writer, "\t%s\t%d\n", name, resp.Namespaces[name].VectorCount)
		}
	}

	writer.Flush()
}
