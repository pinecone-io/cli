package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
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
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Name\t%s\n", ns.Name)
	fmt.Fprintf(writer, "Record Count\t%d\n", ns.RecordCount)

	indexedFieldsVal := "<none>"
	if ns.IndexedFields != nil {
		indexedFieldsVal = text.InlineJSON(ns.IndexedFields)
	}
	fmt.Fprintf(writer, "Indexed Fields\t%s\n", indexedFieldsVal)

	schemaVal := "<none>"
	if ns.Schema != nil {
		schemaVal = text.InlineJSON(ns.Schema)
	}
	fmt.Fprintf(writer, "Schema\t%s\n", schemaVal)

	writer.Flush()
}
