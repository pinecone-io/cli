package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

func PrintDescribeCollectionTable(coll *pinecone.Collection) {
	writer := NewTabWriter()
	log.Debug().Str("name", coll.Name).Msg("Printing collection description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", coll.Name)
	pcio.Fprintf(writer, "State\t%s\n", ColorizeCollectionStatus(coll.Status))
	pcio.Fprintf(writer, "Dimension\t%d\n", coll.Dimension)
	pcio.Fprintf(writer, "Size\t%d\n", coll.Size)
	pcio.Fprintf(writer, "VectorCount\t%d\n", coll.VectorCount)
	pcio.Fprintf(writer, "Environment\t%s\n", coll.Environment)

	writer.Flush()
}

func ColorizeCollectionStatus(state pinecone.CollectionStatus) string {
	switch state {
	case pinecone.CollectionStatusReady:
		return style.SuccessStyle().Render(string(state))
	case pinecone.CollectionStatusInitializing, pinecone.CollectionStatusTerminating:
		return style.WarningStyle().Render(string(state))
	}

	return string(state)
}
