package collection

import (
	"context"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/pinecone"
)

type ListCollectionsCmdOptions struct {
	json bool
}

func NewListCollectionsCmd() *cobra.Command {
	options := ListCollectionsCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "See the list of collections in your project",
		Run: func(cmd *cobra.Command, args []string) {
			pc := sdk.NewPineconeClient()
			ctx := context.Background()

			collections, err := pc.ListCollections(ctx)
			if err != nil {
				exit.Error(err)
			}

			// Sort results alphabetically by name
			sort.SliceStable(collections, func(i, j int) bool {
				return collections[i].Name < collections[j].Name
			})

			if options.json {
				text.PrettyPrintJSON(collections)
			} else {
				printTable(collections)
			}
		},
	}

	// Optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func int32ToString(i *int32) string {
	if i == nil {
		return ""
	}
	return strconv.FormatInt(int64(*i), 10)
}

func int64ToString(i *int64) string {
	if i == nil {
		return ""
	}
	return strconv.FormatInt(*i, 10)
}

func printTable(collections []*pinecone.Collection) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"NAME", "DIMENSION", "SIZE", "STATUS", "VECTORS", "ENVIRONMENT"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, coll := range collections {
		values := []string{coll.Name, int32ToString(coll.Dimension), int64ToString(coll.Size), string(coll.Status), int32ToString(coll.VectorCount), coll.Environment}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}
