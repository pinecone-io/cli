package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
)

func PrintDescribeAssistantTable(am *assistants.AssistantModel) {

	writer := NewTabWriter()

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", style.Emphasis(am.Name))
	pcio.Fprintf(writer, "Metadata\t%s\n", text.InlineJSON(am.Metadata))
	pcio.Fprintf(writer, "Status\t%s\n", am.Status)
	pcio.Fprintf(writer, "\t\n")
	pcio.Fprintf(writer, "CreatedAt\t%s\n", am.CreatedAt)
	pcio.Fprintf(writer, "UpdatedAt\t%s\n", am.UpdatedAt)
	pcio.Fprintf(writer, "\t\n")

	writer.Flush()
}
