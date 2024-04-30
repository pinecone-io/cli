package index

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type DeleteCmdOptions struct {
	name string
}

func NewDeleteCmd() *cobra.Command {
	options := DeleteCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an index",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			err := pc.DeleteIndex(ctx, options.name)
			if err != nil {
				exit.Error(err)
			}

			pcio.Printf(style.SuccessMsg("Index %s deleted.\n"), style.Emphasis(options.name))
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to describe")
	cmd.MarkFlagRequired("name")

	return cmd
}
