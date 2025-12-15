package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintDescribeNamespaceTable(ns *pinecone.NamespaceDescription) {
	writer := NewTabWriter()
	if ns == nil {
		PrintEmptyState(writer, "namespace details")
		return
	}

	log.Debug().Str("name", ns.Name).Msg("Printing namespace description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", ns.Name)
	pcio.Fprintf(writer, "Record Count\t%d\n", ns.RecordCount)

	indexedFieldsVal := "<none>"
	if ns.IndexedFields != nil {
		indexedFieldsVal = text.InlineJSON(ns.IndexedFields)
	}
	pcio.Fprintf(writer, "Indexed Fields\t%s\n", indexedFieldsVal)

	schemaVal := "<none>"
	if ns.Schema != nil {
		schemaVal = text.InlineJSON(ns.Schema)
	}
	pcio.Fprintf(writer, "Schema\t%s\n", schemaVal)

	writer.Flush()
}
