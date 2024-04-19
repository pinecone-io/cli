package describe

import (
	"fmt"
	"os"
	"context"

	"github.com/spf13/cobra"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

var helpText = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

var indexName string

func NewDescribeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Get configuration and status information for an index",
		Long: helpText,
		Run: func(cmd *cobra.Command, args []string) {
			key := os.Getenv("PINECONE_API_KEY")
			fmt.Println("describe called with key:", key)
			fmt.Println("describe called with index name:", indexName)

			ctx := context.Background()

			pc, err := pinecone.NewClient(pinecone.NewClientParams{
				ApiKey: key,
			})
		
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		
			idxs, err := pc.ListIndexes(ctx)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(idxs)
		
			// for _, index := range idxs {
			// 	fmt.Println(index)
			// }
		
			// idx, err := pc.Index(idxs[0].Host)
			// defer idx.Close()
		
			// if err != nil {
			// 	fmt.Println("Error:", err)
			// 	return
			// }
		
			// res, err := idx.DescribeIndexStats(&ctx)
			// if err != nil {
			// 	fmt.Println("Error:", err)
			// 	return
			// }
		
			// fmt.Println(res)
		},
	}
	cmd.Flags().StringVarP(&indexName, "name", "n", "", "name of index to describe")
	cmd.MarkFlagRequired("name")

	return cmd
}