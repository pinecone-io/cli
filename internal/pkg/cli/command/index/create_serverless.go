package index

import (
	"context"
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/client"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/spf13/cobra"
)

type describeOptions struct {
	name      string
	dimension int32
	metric    string
	cloud     string
	region    string
	json      bool
}

func (o describeOptions) isValidCloud() bool {
	switch pinecone.Cloud(o.cloud) {
	case pinecone.Aws, pinecone.Azure, pinecone.Gcp:
		return true
	default:
		return false
	}
}

func (o describeOptions) isValidMetric() bool {
	switch pinecone.IndexMetric(o.metric) {
	case pinecone.Cosine, pinecone.Euclidean, pinecone.Dotproduct:
		return true
	default:
		return false
	}
}

func NewCreateServerlessCmd() *cobra.Command {
	options := describeOptions{}

	cmd := &cobra.Command{
		Use:   "create-serverless",
		Short: "Create a serverless index with the specified configuration",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !options.isValidCloud() {
				return fmt.Errorf("cloud provider must be one of [aws, azure, gcp]")
			}
			if !options.isValidMetric() {
				return fmt.Errorf("metric must be one of [cosine, euclidean, dotproduct]")
			}

			return nil
		},
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
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func runCreateServerlessCmd(cmd *cobra.Command, options describeOptions) {
	ctx := context.Background()
	pc := client.NewPineconeClient()

	createRequest := &pinecone.CreateServerlessIndexRequest{
		Name:      options.name,
		Metric:    pinecone.IndexMetric(options.metric),
		Dimension: options.dimension,
		Cloud:     pinecone.Cloud(options.cloud),
		Region:    options.region,
	}

	idx, err := pc.CreateServerlessIndex(ctx, createRequest)
	if err != nil {
		exit.Error(err)
	}
	if options.json {
		text.PrettyPrintJSON(idx)
		return
	}

	describeCommand := fmt.Sprintf("pinecone index describe --name %s", idx.Name)
	fmt.Fprintf(cmd.OutOrStdout(), "✅ Index %s created successfully. Run %s to monitor status. \n\n", style.Emphasis(idx.Name), style.Code(describeCommand))
	presenters.PrintDescribeIndexTable(idx)
}
