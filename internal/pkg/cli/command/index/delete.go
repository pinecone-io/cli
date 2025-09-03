package index

import (
	"context"
	"fmt"

	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/index"
	"github.com/pinecone-io/cli/internal/pkg/utils/interactive"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
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
		Use:          "delete <name>",
		Short:        "Delete an index",
		Args:         index.ValidateIndexNameArgs,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			options.name = args[0]

			// Ask for user confirmation
			question := fmt.Sprintf("Do you want to delete the index '%s'?", options.name)
			if !interactive.GetConfirmation(question) {
				pcio.Println(style.InfoMsg("Index deletion cancelled."))
				return
			}

			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			err := pc.DeleteIndex(ctx, options.name)
			if err != nil {
				errorutil.HandleIndexAPIError(err, cmd, args)
				exit.Error(err)
			}

			msg.SuccessMsg("Index %s deleted.\n", style.Emphasis(options.name))
		},
	}

	return cmd
}
