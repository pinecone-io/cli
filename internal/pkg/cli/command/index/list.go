package index

import (
	"context"
	"sort"

	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
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
				errorutil.HandleIndexAPIError(err, cmd, []string{})
				exit.Error(err)
			}

			// Sort results alphabetically by name
			sort.SliceStable(idxs, func(i, j int) bool {
				return idxs[i].Name < idxs[j].Name
			})

			if options.json {
				json := text.IndentJSON(idxs)
				pcio.Println(json)
			} else {
				// Show essential and state information
				presenters.PrintIndexTableWithIndexAttributesGroups(idxs, []presenters.IndexAttributesGroup{
					presenters.IndexAttributesGroupEssential,
					presenters.IndexAttributesGroupState,
				})
			}
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
