package index

import (
	"context"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/pinecone"
)

type ListIndexCmdOptions struct {
	json bool
}

func NewListCmd() *cobra.Command {
	options := ListIndexCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "See the list of indexes in the targeted project",
		Run: func(cmd *cobra.Command, args []string) {
			pc := sdk.NewPineconeClient()
			ctx := context.Background()

			idxs, err := pc.ListIndexes(ctx)
			if err != nil {
				msg.FailMsg("Failed to list indexes: %s\n", err)
				exit.Error(err)
			}

			// Sort results alphabetically by name
			sort.SliceStable(idxs, func(i, j int) bool {
				return idxs[i].Name < idxs[j].Name
			})

			if options.json {
				text.PrettyPrintJSON(idxs)
			} else {
				printTable(idxs)
			}
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func printTable(idxs []*pinecone.Index) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"NAME", "STATUS", "HOST", "DIMENSION", "METRIC", "SPEC"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, idx := range idxs {
		if idx.Spec.Serverless == nil {
			// Pod index
			values := []string{idx.Name, string(idx.Status.State), idx.Host, pcio.Sprintf("%d", idx.Dimension), string(idx.Metric), "pod"}
			pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
		} else {
			// Serverless index
			values := []string{idx.Name, string(idx.Status.State), idx.Host, pcio.Sprintf("%d", idx.Dimension), string(idx.Metric), "serverless"}
			pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
		}
	}
	writer.Flush()
}
