package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintUpdateVectorsByMetadataTable(resp *pinecone.UpdateVectorsByMetadataResponse) {
	writer := NewTabWriter()
	if resp == nil {
		PrintEmptyState(writer, "update results")
		return
	}

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Matched Records\t%d\n", resp.MatchedRecords)

	writer.Flush()
}
