package collection

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
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
		Short: "Describe a collection by name",
		Example: help.Examples(`
			pc collection describe --name "collection-name"
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			collection, err := pc.DescribeCollection(ctx, options.name)
			if err != nil {
				msg.FailMsg("Failed to describe collection %s: %s\n", options.name, err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(collection)
				pcio.Println(json)
			} else {
				presenters.PrintDescribeCollectionTable(collection)
			}
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of collection to describe")
	_ = cmd.MarkFlagRequired("name")

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
