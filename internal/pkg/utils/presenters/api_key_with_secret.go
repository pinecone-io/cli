package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

func PrintDescribeAPIKeyWithSecretTable(apiKey *pinecone.APIKeyWithSecret) {
	writer := NewTabWriter()
	log.Debug().Str("name", apiKey.Key.Name).Msg("Printing API key description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", apiKey.Key.Name)
	pcio.Fprintf(writer, "ID\t%s\n", apiKey.Key.Id)
	pcio.Fprintf(writer, "Value\t%s\n", apiKey.Value)
	pcio.Fprintf(writer, "Project ID\t%s\n", apiKey.Key.ProjectId)
	pcio.Fprintf(writer, "Roles\t%s\n", strings.Join(apiKey.Key.Roles, ", "))

	writer.Flush()
}
