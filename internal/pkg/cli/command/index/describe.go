package index

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeCmdOptions struct {
	json bool
}

func NewDescribeCmd() *cobra.Command {
	options := DescribeCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe [name]",
		Short: "Get configuration and status information for an index",
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
			idx, err := pc.DescribeIndex(ctx, name)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					msg.FailMsg("The index %s does not exist\n", style.Emphasis(name))
				} else {
					msg.FailMsg("Failed to describe index %s: %s\n", style.Emphasis(name), err)
				}
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(idx)
				pcio.Println(json)
			} else {
				presenters.PrintDescribeIndexTable(idx)
			}
		},
	}

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
