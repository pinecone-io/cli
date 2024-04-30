package index

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/spf13/cobra"
)

type createServerlessOptions struct {
	name      string
	dimension int32
	metric    string
	cloud     string
	region    string
	json      bool
}

func NewCreateServerlessCmd() *cobra.Command {
	options := createServerlessOptions{}

	cmd := &cobra.Command{
		Use:   "create-serverless",
		Short: "Create a serverless index with the specified configuration",
		Run: func(cmd *cobra.Command, args []string) {
			runCreateServerlessCmd(cmd, options)
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the index")
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

func runCreateServerlessCmd(cmd *cobra.Command, options createServerlessOptions) {
	ctx := context.Background()
	pc := sdk.NewPineconeClient()

	createRequest := &pinecone.CreateServerlessIndexRequest{
		Name:      options.name,
		Metric:    pinecone.IndexMetric(options.metric),
		Dimension: options.dimension,
		Cloud:     pinecone.Cloud(options.cloud),
		Region:    options.region,
	}

	idx, err := pc.CreateServerlessIndex(ctx, createRequest)
	if err != nil {
		msg.FailMsg("Failed to create serverless index %s: %s\n", style.Emphasis(options.name), err)
		exit.Error(err)
	}
	if options.json {
		text.PrettyPrintJSON(idx)
		return
	}

	describeCommand := pcio.Sprintf("pinecone index describe --name %s", idx.Name)
	msg.SuccessMsg("Index %s created successfully. Run %s to check status. \n\n", style.Emphasis(idx.Name), style.Code(describeCommand))
	presenters.PrintDescribeIndexTable(idx)
}
