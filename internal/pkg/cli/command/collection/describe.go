package collection

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/client"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeCollectionCmdOptions struct {
	name string
	json bool
}

func NewDescribeCollectionCmd() *cobra.Command {
	options := DescribeCollectionCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Get information on a collection",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			pc := client.NewPineconeClient()

			collection, err := pc.DescribeCollection(ctx, options.name)
			if err != nil {
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(collection)
			} else {
				presenters.PrintDescribeCollectionTable(collection)
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of collection to describe")
	cmd.MarkFlagRequired("name")

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
