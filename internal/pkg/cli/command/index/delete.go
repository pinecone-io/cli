package index

import (
	"context"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
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
				if strings.Contains(err.Error(), "not found") {
					msg.FailMsg("The index %s does not exist\n", style.Emphasis(options.name))
				} else {
					msg.FailMsg("Failed to delete index %s: %s\n", style.Emphasis(options.name), err)
				}
				exit.Error(err)
			}

			msg.SuccessMsg("Index %s deleted.\n", style.Emphasis(options.name))
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to delete")
	cmd.MarkFlagRequired("name")

	return cmd
}
