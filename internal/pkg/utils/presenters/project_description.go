package presenters

import (
	"strconv"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

func PrintDescribeProjectTable(proj *pinecone.Project) {
	writer := NewTabWriter()
	log.Debug().Str("name", proj.Name).Msg("Printing project description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", proj.Name)
	pcio.Fprintf(writer, "ID\t%s\n", proj.Id)
	pcio.Fprintf(writer, "Organization ID\t%s\n", proj.OrganizationId)
	pcio.Fprintf(writer, "Created At\t%s\n", proj.CreatedAt.String())
	pcio.Fprintf(writer, "Force Encryption\t%s\n", strconv.FormatBool(proj.ForceEncryptionWithCmek))
	pcio.Fprintf(writer, "Max Pods\t%d\n", proj.MaxPods)
	pcio.Fprintf(writer, "\t\n")

	writer.Flush()
}
