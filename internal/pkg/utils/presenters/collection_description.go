package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintDescribeCollectionTable(coll *pinecone.Collection) {
	writer := NewTabWriter()
	if coll == nil {
		PrintEmptyState(writer, "collection details")
		return
	}

	log.Debug().Str("name", coll.Name).Msg("Printing collection description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Name\t%s\n", coll.Name)
	fmt.Fprintf(writer, "State\t%s\n", ColorizeCollectionStatus(coll.Status))
	fmt.Fprintf(writer, "Dimension\t%d\n", coll.Dimension)
	fmt.Fprintf(writer, "Size\t%d\n", coll.Size)
	fmt.Fprintf(writer, "VectorCount\t%d\n", coll.VectorCount)
	fmt.Fprintf(writer, "Environment\t%s\n", coll.Environment)

	writer.Flush()
}

func ColorizeCollectionStatus(state pinecone.CollectionStatus) string {
	switch state {
	case pinecone.CollectionStatusReady:
		return style.StatusGreen(string(state))
	case pinecone.CollectionStatusInitializing, pinecone.CollectionStatusTerminating:
		return style.StatusYellow(string(state))
	}

	return string(state)
}
