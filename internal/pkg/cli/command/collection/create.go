package collection

import (
	"context"
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/collection/presenters"
	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/v4/pinecone"
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
				errorutil.HandleAPIError(err, cmd, args)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(collection)
				pcio.Println(json)
			} else {
				renderSuccessOutput(collection)
			}
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name you want to give the collection")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&options.sourceIndex, "source", "s", "", "name of the pod-based index to use as the source for the collection")
	_ = cmd.MarkFlagRequired("source")

	// Optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func renderSuccessOutput(collection *pinecone.Collection) {
	msg.SuccessMsg("Collection %s created successfully.", style.ResourceName(collection.Name))

	presenters.PrintDescribeCollectionTable(collection)

	describeCommand := pcio.Sprintf("pc collection describe %s", collection.Name)
	hint := fmt.Sprintf("Run %s at any time to check the status. \n\n", style.Code(describeCommand))
	pcio.Println(style.Hint(hint))
}
