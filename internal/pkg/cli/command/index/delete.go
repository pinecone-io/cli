package index

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete an index",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("index name is required")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			name := args[0]
			err := pc.DeleteIndex(ctx, name)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					msg.FailMsg("The index %s does not exist\n", style.Emphasis(name))
				} else {
					msg.FailMsg("Failed to delete index %s: %s\n", style.Emphasis(name), err)
				}
				exit.Error(err)
			}

			msg.SuccessMsg("Index %s deleted.\n", style.Emphasis(name))
		},
	}

	return cmd
}
