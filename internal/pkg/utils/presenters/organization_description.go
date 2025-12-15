package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func PrintDescribeOrganizationTable(org *pinecone.Organization) {
	writer := NewTabWriter()
	if org == nil {
		PrintEmptyState(writer, "organization details")
		return
	}

	log.Debug().Str("name", org.Name).Msg("Printing organization description")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", org.Name)
	pcio.Fprintf(writer, "ID\t%s\n", org.Id)
	pcio.Fprintf(writer, "Created At\t%s\n", org.CreatedAt.String())
	pcio.Fprintf(writer, "Payment Status\t%s\n", org.PaymentStatus)
	pcio.Fprintf(writer, "Plan\t%s\n", org.Plan)
	pcio.Fprintf(writer, "Support Tier\t%s\n", org.SupportTier)
	pcio.Fprintf(writer, "\n")

	writer.Flush()

}
