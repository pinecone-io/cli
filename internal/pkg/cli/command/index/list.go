package index

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/utils/client"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/pinecone"
)

var listHelpText = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

type ListIndexCmdOptions struct {
	json bool
}

func NewListCmd() *cobra.Command {
	options := ListIndexCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "See the list of indexes in your project",
		Long:  listHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			pc := client.NewPineconeClient()
			ctx := context.Background()

			idxs, err := pc.ListIndexes(ctx)
			if err != nil {
				exit.Error(err)
			}

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

	columns := []string{"NAME", "STATUS", "HOST", "DIMENSION", "METRIC", "SPEC", "CLOUD", "REGION", "ENVIRONMENT"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	for _, idx := range idxs {
		if idx.Spec.Serverless == nil {
			// Pod index
			values := []string{idx.Name, string(idx.Status.State), idx.Host, fmt.Sprintf("%d", idx.Dimension), string(idx.Metric), "pod", "", "", idx.Spec.Pod.Environment}
			fmt.Fprintf(writer, strings.Join(values, "\t")+"\n")
		} else {
			// Serverless index
			values := []string{idx.Name, string(idx.Status.State), idx.Host, fmt.Sprintf("%d", idx.Dimension), string(idx.Metric), "serverless", string(idx.Spec.Serverless.Cloud), idx.Spec.Serverless.Region, ""}
			fmt.Fprintf(writer, strings.Join(values, "\t")+"\n")
		}
	}
	writer.Flush()
}
