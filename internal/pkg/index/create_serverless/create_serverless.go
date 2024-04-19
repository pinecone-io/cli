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
var dimension int32
var metric string
var cloud string
var region string

func NewCreateServerlessCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-serverless",
		Short: "Create a serverless index with the specified configuration",
		Long: helpText,
		Run: func(cmd *cobra.Command, args []string) {
			key := os.Getenv("PINECONE_API_KEY")
			fmt.Println("describe called with key:", key)
			fmt.Println("describe called with index name:", indexName)
			fmt.Println("describe called with dimension:", dimension)
			fmt.Println("describe called with metric:", metric)
			fmt.Println("describe called with cloud:", cloud)
			fmt.Println("describe called with region:", region)

			ctx := context.Background()

			pc, err := pinecone.NewClient(pinecone.NewClientParams{
				ApiKey: key,
			})
		
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			createRequest := &pinecone.CreateServerlessIndexRequest{
				Name: indexName,
				Metric: pinecone.Cosine,
				Dimension: dimension,
				Cloud: pinecone.Aws,
				Region: region,
			}
			idx, err := pc.CreateServerlessIndex(ctx, createRequest)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(idx)
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&indexName, "name", "n", "", "name of index to describe")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&cloud, "cloud", "c", "", "cloud provider where you would like to deploy")
	cmd.MarkFlagRequired("cloud")
	cmd.Flags().StringVarP(&region, "region", "r", "", "cloud region where you would like to deploy")
	cmd.MarkFlagRequired("region")
	cmd.Flags().Int32VarP(&dimension, "dimension", "d", 0, "dimension of the index to create")
	cmd.MarkFlagRequired("dimension")

	// Optional flags
	cmd.Flags().StringVarP(&metric, "metric", "m", "cosine", "metric to use. One of: cosine, euclidean, dotproduct")

	return cmd
}