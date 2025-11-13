package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintDescribeAPIKeyTable(apiKey *pinecone.APIKey) {
	writer := NewTabWriter()
	log.Debug().Str("name", apiKey.Name).Msg("Printing API key description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", apiKey.Name)
	pcio.Fprintf(writer, "ID\t%s\n", apiKey.Id)
	pcio.Fprintf(writer, "Project ID\t%s\n", apiKey.ProjectId)
	pcio.Fprintf(writer, "Roles\t%s\n", strings.Join(apiKey.Roles, ", "))

	writer.Flush()
}
