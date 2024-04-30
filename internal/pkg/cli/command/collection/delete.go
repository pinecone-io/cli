package collection

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type DeleteCollectionCmdOptions struct {
	name string
	json bool
}

func NewDeleteCollectionCmd() *cobra.Command {
	options := DeleteCollectionCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a collection",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			err := pc.DeleteCollection(ctx, options.name)
			if err != nil {
				exit.Error(err)
			}

			pcio.Printf(style.SuccessMsg("Collection %s deleted.\n"), style.Emphasis(options.name))
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of collection to delete")
	cmd.MarkFlagRequired("name")

	return cmd
}
