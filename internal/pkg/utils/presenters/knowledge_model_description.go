package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
)

func PrintDescribeKnowledgeModelTable(km *knowledge.KnowledgeModel) {

	writer := NewTabWriter()
	log.Debug().Str("name", km.Name).Msg("Printing knowledge model description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", style.Emphasis(km.Name))
	pcio.Fprintf(writer, "Metadata\t%s\n", text.InlineJSON(km.Metadata))
	pcio.Fprintf(writer, "Status\t%s\n", km.Status)
	pcio.Fprintf(writer, "\t\n")
	pcio.Fprintf(writer, "CreatedAt\t%s\n", km.CreatedAt)
	pcio.Fprintf(writer, "UpdatedAt\t%s\n", km.UpdatedAt)
	pcio.Fprintf(writer, "\t\n")

	writer.Flush()
}
