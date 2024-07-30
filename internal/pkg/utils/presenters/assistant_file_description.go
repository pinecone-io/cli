package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

func PrintDescribeAssistantFileTable(file *assistants.AssistantFileModel) {
	writer := NewTabWriter()

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", file.Name)
	pcio.Fprintf(writer, "Id\t%s\n", file.Id)
	pcio.Fprintf(writer, "Metadata\t%s\n", file.Metadata.ToString())
	pcio.Fprintf(writer, "CreatedOn\t%s\n", file.CreatedOn)
	pcio.Fprintf(writer, "UpdatedOn\t%s\n", file.UpdatedOn)
	pcio.Fprintf(writer, "Status\t%s\n", file.Status)
	pcio.Fprintf(writer, "Size\t%d\n", file.Size)

	writer.Flush()
}
