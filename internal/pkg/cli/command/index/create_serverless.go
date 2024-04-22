package index

import (
	"fmt"
	"os"
	"context"

	"github.com/spf13/cobra"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
)

var serverlessHelpText = `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

type describeOptions struct {
	name string
	dimension int32
	metric string
	cloud string
	region string
}

func NewCreateServerlessCmd() *cobra.Command {
	options := describeOptions{}

	cmd := &cobra.Command{
		Use:   "create-serverless",
		Short: "Create a serverless index with the specified configuration",
		Long: serverlessHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			runCreateServerlessCmd(cmd, options)
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to describe")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&options.cloud, "cloud", "c", "", "cloud provider where you would like to deploy")
	cmd.MarkFlagRequired("cloud")
	cmd.Flags().StringVarP(&options.region, "region", "r", "", "cloud region where you would like to deploy")
	cmd.MarkFlagRequired("region")
	cmd.Flags().Int32VarP(&options.dimension, "dimension", "d", 0, "dimension of the index to create")
	cmd.MarkFlagRequired("dimension")

	// Optional flags
	cmd.Flags().StringVarP(&options.metric, "metric", "m", "cosine", "metric to use. One of: cosine, euclidean, dotproduct")

	return cmd
}

func runCreateServerlessCmd(cmd *cobra.Command, options describeOptions) {
	key := os.Getenv("PINECONE_API_KEY")
	fmt.Println("describe called with key:", key)
	fmt.Println("describe called with index name:", options)
	// fmt.Println("describe called with dimension:", dimension)
	// fmt.Println("describe called with metric:", metric)
	// fmt.Println("describe called with cloud:", cloud)
	// fmt.Println("describe called with region:", region)

	ctx := context.Background()

	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: key,
	})

	if err != nil {
		exit.Error(err)
	}

	createRequest := &pinecone.CreateServerlessIndexRequest{
		Name: options.name,
		Metric: pinecone.Cosine,
		Dimension: options.dimension,
		Cloud: pinecone.Aws,
		Region: options.region,
	}

	idx, err := pc.CreateServerlessIndex(ctx, createRequest)
	if err != nil {
		exit.Error(err)
	}
	fmt.Println(idx)
}