package collection

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/pinecone"
)

type CreateCollectionCmdOptions struct {
	json        bool
	name        string
	sourceIndex string
}

func NewCreateCollectionCmd() *cobra.Command {
	options := CreateCollectionCmdOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a collection from a pod-based index",
		Run: func(cmd *cobra.Command, args []string) {
			pc := sdk.NewPineconeClient()
			ctx := context.Background()

			req := &pinecone.CreateCollectionRequest{
				Name:   options.name,
				Source: options.sourceIndex,
			}
			collection, err := pc.CreateCollection(ctx, req)
			if err != nil {
				msg.FailMsg("Failed to create collection: %s\n", err)
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(collection)
			} else {
				describeCommand := pcio.Sprintf("pinecone collection describe --name %s", collection.Name)
				msg.SuccessMsg("Collection %s created successfully. Run %s to check status. \n\n", style.Emphasis(collection.Name), style.Code(describeCommand))
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
