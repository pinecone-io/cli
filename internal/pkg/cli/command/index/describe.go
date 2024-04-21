package index

import (
	"fmt"
	"os"
	"context"

	"github.com/spf13/cobra"
	"github.com/pinecone-io/go-pinecone/pinecone"
	text "github.com/pinecone-io/cli/internal/pkg/utils/text"
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
			key := os.Getenv("PINECONE_API_KEY")
			fmt.Println("describe called with key:", key)
			fmt.Println("describe called with index name:", options.name)

			ctx := context.Background()

			pc, err := pinecone.NewClient(pinecone.NewClientParams{
				ApiKey: key,
			})
		
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		
			idxs, err := pc.DescribeIndex(ctx, options.name)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			text.PrettyPrintJSON(idxs)
		},
	}
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to describe")
	cmd.MarkFlagRequired("name")

	return cmd
}