package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintDescribeAPIKeyTable(apiKey *pinecone.APIKey) {
	writer := NewTabWriter()
	if apiKey == nil {
		PrintEmptyState(writer, "API key details")
		return
	}

	log.Debug().Str("name", apiKey.Name).Msg("Printing API key description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Name\t%s\n", apiKey.Name)
	fmt.Fprintf(writer, "ID\t%s\n", apiKey.Id)
	fmt.Fprintf(writer, "Project ID\t%s\n", apiKey.ProjectId)
	fmt.Fprintf(writer, "Roles\t%s\n", strings.Join(apiKey.Roles, ", "))

	writer.Flush()
}
