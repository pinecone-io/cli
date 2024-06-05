package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

func PrintDescribeKnowledgeFileTable(file *knowledge.KnowledgeFileModel) {
	writer := NewTabWriter()
	log.Debug().Str("file_id", file.Id).Msg("Printing knowledge file description")

	columns := []string{"ATTIRBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", file.Name)
	pcio.Fprintf(writer, "Id\t%s\n", file.Id)
	pcio.Fprintf(writer, "Metadata\t%s\n", file.Metadata.ToString())
	pcio.Fprintf(writer, "CreatedOn\t%s\n", file.CreatedOn)
	pcio.Fprintf(writer, "UpdatedOn\t%s\n", file.UpdatedOn)
	pcio.Fprintf(writer, "Status\t%s\n", file.Status)

	writer.Flush()
}
