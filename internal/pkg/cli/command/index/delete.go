package index

import (
	"context"
	"errors"
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
		Use:   "delete <name>",
		Short: "Delete an index",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				// TODO: start interactive mode. For now just return an error.
				return errors.New("please provide an index name")
			}
			if len(args) > 1 {
				return errors.New("please provide only one index name")
			}
			if strings.TrimSpace(args[0]) == "" {
				return errors.New("index name cannot be empty")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			options.name = args[0]
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

	return cmd
}
