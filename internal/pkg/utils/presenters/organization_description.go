package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
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
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Name\t%s\n", org.Name)
	fmt.Fprintf(writer, "ID\t%s\n", org.Id)
	fmt.Fprintf(writer, "Created At\t%s\n", org.CreatedAt.String())
	fmt.Fprintf(writer, "Payment Status\t%s\n", org.PaymentStatus)
	fmt.Fprintf(writer, "Plan\t%s\n", org.Plan)
	fmt.Fprintf(writer, "Support Tier\t%s\n", org.SupportTier)
	fmt.Fprintf(writer, "\n")

	writer.Flush()

}
