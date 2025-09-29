package collection

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/collection"
	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewDeleteCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "delete <name>",
		Short:        "Delete a collection",
		Args:         collection.ValidateCollectionNameArgs,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			collectionName := args[0]
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			err := pc.DeleteCollection(ctx, collectionName)
			if err != nil {
				errorutil.HandleAPIError(err, cmd, args)
				exit.Error(err)
			}

			msg.SuccessMsg("Collection %s deleted.\n", style.Emphasis(collectionName))
		},
	}

	return cmd
}
