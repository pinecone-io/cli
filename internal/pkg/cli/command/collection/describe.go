package collection

import (
	"context"
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/collection"
	"github.com/pinecone-io/cli/internal/pkg/utils/collection/presenters"
	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeCollectionCmdOptions struct {
	json bool
}

func NewDescribeCollectionCmd() *cobra.Command {
	options := DescribeCollectionCmdOptions{}

	cmd := &cobra.Command{
		Use:          "describe <name>",
		Short:        "Get information on a collection",
		Args:         collection.ValidateCollectionNameArgs,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			collectionName := args[0]
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			collection, err := pc.DescribeCollection(ctx, collectionName)
			if err != nil {
				errorutil.HandleAPIError(err, cmd, args)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(collection)
				fmt.Println(json)
			} else {
				presenters.PrintDescribeCollectionTable(collection)
			}
		},
	}

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
