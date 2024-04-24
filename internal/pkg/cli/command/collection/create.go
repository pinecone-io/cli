package collection

import (
	"context"
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/client"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/pinecone"
)

type CreateCollectionOptions struct {
	json        bool
	name        string
	sourceIndex string
}

func NewCreateCollectionCmd() *cobra.Command {
	options := CreateCollectionOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a collection from a pod-based index",
		Run: func(cmd *cobra.Command, args []string) {
			pc := client.NewPineconeClient()
			ctx := context.Background()

			req := &pinecone.CreateCollectionRequest{
				Name:   options.name,
				Source: options.sourceIndex,
			}
			collection, err := pc.CreateCollection(ctx, req)
			if err != nil {
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(collection)
			} else {
				describeCommand := fmt.Sprintf("pinecone collection describe --name %s", collection.Name)
				fmt.Fprintf(cmd.OutOrStdout(), "✅ Collection %s created successfully. Run %s to monitor status. \n\n", style.Emphasis(collection.Name), style.Code(describeCommand))
				presenters.PrintDescribeCollectionTable(collection)
			}
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name you want to give the collection")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&options.sourceIndex, "source", "s", "", "name of the index to use as the source for the collection")
	cmd.MarkFlagRequired("source")

	// Optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}