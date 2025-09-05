package index

import (
	"fmt"

	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/index"
	indexpresenters "github.com/pinecone-io/cli/internal/pkg/utils/index/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
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
		Use:          "describe <name>",
		Short:        "Get configuration and status information for an index",
		Args:         index.ValidateIndexNameArgs,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			options.name = args[0]
			pc := sdk.NewPineconeClient()

			idx, err := pc.DescribeIndex(cmd.Context(), options.name)
			if err != nil {
				errorutil.HandleIndexAPIError(err, cmd, args)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(idx)
				fmt.Println(json)
			} else {
				indexpresenters.PrintDescribeIndexTable(idx)
			}
		},
	}

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
