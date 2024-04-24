package collection

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/utils/client"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
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
			pc := client.NewPineconeClient()
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

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func printTable(collections []*pinecone.Collection) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"NAME", "DIMENSION", "SIZE", "STATUS", "VECTORS", "ENVIRONMENT"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	for _, coll := range collections {
		values := []string{coll.Name, strconv.FormatInt(int64(*coll.Dimension), 10), strconv.FormatInt(*coll.Size, 10), strconv.FormatInt(int64(*coll.VectorCount), 10), coll.Environment}
		fmt.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}