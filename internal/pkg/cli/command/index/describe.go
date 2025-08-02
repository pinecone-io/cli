package index

import (
	"context"
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
	name string
	json bool
}

func NewDescribeCmd() *cobra.Command {
	options := DescribeCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Get configuration and status information for an index",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			pc := sdk.NewPineconeClient()

			idx, err := pc.DescribeIndex(ctx, options.name)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					msg.FailMsg("The index %s does not exist\n", style.Emphasis(options.name))
				} else {
					msg.FailMsg("Failed to describe index %s: %s\n", style.Emphasis(options.name), err)
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

	// required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to describe")
	_ = cmd.MarkFlagRequired("name")

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
