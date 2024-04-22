package index

import (
	"context"

	"github.com/spf13/cobra"
	text "github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/client"
)

var describeHelpText = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

type DescribeCmdOptions struct {
	name string
}

func NewDescribeCmd() *cobra.Command {
	options := DescribeCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Get configuration and status information for an index",
		Long: describeHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			pc := client.NewPineconeClient()

			idxs, err := pc.DescribeIndex(ctx, options.name)
			if err != nil {
				exit.Error(err)
			}
			text.PrettyPrintJSON(idxs)
		},
	}
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to describe")
	cmd.MarkFlagRequired("name")

	return cmd
}