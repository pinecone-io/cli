package index

import (
	"context"
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/client"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
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
			pc := client.NewPineconeClient()

			err := pc.DeleteIndex(ctx, options.name)
			if err != nil {
				exit.Error(err)
			}

			fmt.Printf("Index %s deleted\n", options.name)
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to describe")
	cmd.MarkFlagRequired("name")

	return cmd
}
