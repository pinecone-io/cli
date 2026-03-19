package presenters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintDescribeProjectTable(proj *pinecone.Project) {
	writer := NewTabWriter()
	if proj == nil {
		PrintEmptyState(writer, "project details")
		return
	}

	log.Debug().Str("name", proj.Name).Msg("Printing project description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Name\t%s\n", proj.Name)
	fmt.Fprintf(writer, "ID\t%s\n", proj.Id)
	fmt.Fprintf(writer, "Organization ID\t%s\n", proj.OrganizationId)
	fmt.Fprintf(writer, "Created At\t%s\n", proj.CreatedAt.String())
	fmt.Fprintf(writer, "Force Encryption\t%s\n", strconv.FormatBool(proj.ForceEncryptionWithCmek))
	fmt.Fprintf(writer, "Max Pods\t%d\n", proj.MaxPods)
	fmt.Fprintf(writer, "\n")

	writer.Flush()
}
