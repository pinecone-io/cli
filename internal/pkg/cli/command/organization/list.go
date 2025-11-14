package organization

import (
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

type listOrganizationCmdOptions struct {
	json bool
}

func NewListOrganizationsCmd() *cobra.Command {
	options := listOrganizationCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all organizations available to the authenticated user",
		Example: help.Examples(`
			pc organization list
		`),
		GroupID: help.GROUP_ORGANIZATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()

			orgs, err := ac.Organization.List(cmd.Context())
			if err != nil {
				msg.FailMsg("Failed to list organizations: %s\n", err)
				exit.Error(err, "Failed to list organizations")
			}

			if options.json {
				json := text.IndentJSON(orgs)
				pcio.Println(json)
				return
			}

			printTable(orgs)
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}

func printTable(orgs []*pinecone.Organization) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"NAME", "ID", "CREATED AT", "PAYMENT STATUS", "PLAN", "SUPPORT TIER"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, org := range orgs {
		values := []string{
			org.Name,
			org.Id,
			org.CreatedAt.String(),
			org.PaymentStatus,
			org.Plan,
			org.SupportTier,
		}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}
