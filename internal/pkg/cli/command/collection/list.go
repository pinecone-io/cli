package collection

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

type listCollectionsCmdOptions struct {
	json bool
}

func NewListCollectionsCmd() *cobra.Command {
	options := listCollectionsCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "See the list of collections in your project",
		Example: help.Examples(`
			pc collection list
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			collections, err := pc.ListCollections(ctx)
			if err != nil {
				msg.FailMsg("Failed to list collections: %s\n", err)
				exit.Error(err, "Failed to list collections")
			}

			// Sort results alphabetically by name
			sort.SliceStable(collections, func(i, j int) bool {
				return collections[i].Name < collections[j].Name
			})

			if options.json {
				json := text.IndentJSON(collections)
				pcio.Println(json)
			} else {
				printTable(collections)
			}
		},
	}

	// Optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func printTable(collections []*pinecone.Collection) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"NAME", "DIMENSION", "SIZE", "STATUS", "VECTORS", "ENVIRONMENT"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, coll := range collections {
		values := []string{coll.Name, string(coll.Dimension), strconv.FormatInt(coll.Size, 10), string(coll.Status), string(coll.VectorCount), coll.Environment}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}
