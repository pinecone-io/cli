package collection

import (
	"context"
	"fmt"
	"sort"

	"github.com/pinecone-io/cli/internal/pkg/utils/collection/presenters"
	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
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
				errorutil.HandleAPIError(err, cmd, args)
				exit.Error(err)
			}

			// Sort results alphabetically by name
			sort.SliceStable(collections, func(i, j int) bool {
				return collections[i].Name < collections[j].Name
			})

			if options.json {
				json := text.IndentJSON(collections)
				fmt.Println(json)
			} else {
				presenters.PrintCollectionTable(collections)
			}
		},
	}

	// Optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
