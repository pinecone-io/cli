package index

import (
	"context"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type deleteCmdOptions struct {
	name string
}

func NewDeleteCmd() *cobra.Command {
	options := deleteCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an index by name",
		Example: help.Examples(`
			pc index delete --name "index-name"
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			err := pc.DeleteIndex(ctx, options.name)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					msg.FailMsg("The index %s does not exist\n", style.Emphasis(options.name))
					exit.Errorf(err, "The index %s does not exist", style.Emphasis(options.name))
				} else {
					msg.FailMsg("Failed to delete index %s: %s\n", style.Emphasis(options.name), err)
					exit.Errorf(err, "Failed to delete index %s", style.Emphasis(options.name))
				}
			}

			msg.SuccessMsg("Index %s deleted.\n", style.Emphasis(options.name))
		},
	}

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to delete")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
