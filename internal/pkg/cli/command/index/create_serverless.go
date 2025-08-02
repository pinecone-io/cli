package index

import (
	"context"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type createServerlessOptions struct {
	name               string
	dimension          int32
	metric             string
	cloud              string
	region             string
	deletionProtection string

	json bool
}

func NewCreateServerlessCmd() *cobra.Command {
	options := createServerlessOptions{}

	cmd := &cobra.Command{
		Use:   "create-serverless",
		Short: "Create a serverless index with the specified configuration",
		Run: func(cmd *cobra.Command, args []string) {
			runCreateServerlessCmd(options)
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the index")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&options.cloud, "cloud", "c", "", "cloud provider where you would like to deploy")
	_ = cmd.MarkFlagRequired("cloud")
	cmd.Flags().StringVarP(&options.region, "region", "r", "", "cloud region where you would like to deploy")
	_ = cmd.MarkFlagRequired("region")
	cmd.Flags().Int32VarP(&options.dimension, "dimension", "d", 0, "dimension of the index to create")
	_ = cmd.MarkFlagRequired("dimension")

	// Optional flags
	cmd.Flags().StringVarP(&options.metric, "metric", "m", "cosine", "metric to use. One of: cosine, euclidean, dotproduct")
	cmd.Flags().StringVarP(&options.deletionProtection, "deletion_protection", "p", "", "Whether to enable deletion protection for the index")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func runCreateServerlessCmd(options createServerlessOptions) {
	ctx := context.Background()
	pc := sdk.NewPineconeClient()

	// Deprecation warning
	pcio.Fprintf(os.Stderr, "⚠️  Warning: The '%s' command is deprecated. Please use '%s' instead.", style.Code("index create-serverless"), style.Code("index create"))

	// Create variables for optional fields that need pointers
	var indexMetric *pinecone.IndexMetric
	if options.metric != "" {
		im := pinecone.IndexMetric(options.metric)
		indexMetric = &im
	}
	var dimension *int32
	if options.dimension != 0 {
		dimension = &options.dimension
	}
	var deletionProtection *pinecone.DeletionProtection
	if options.deletionProtection != "" {
		dp := pinecone.DeletionProtection(options.deletionProtection)
		deletionProtection = &dp
	}

	createRequest := &pinecone.CreateServerlessIndexRequest{
		Name:               options.name,
		Metric:             indexMetric,
		Dimension:          dimension,
		Cloud:              pinecone.Cloud(options.cloud),
		Region:             options.region,
		DeletionProtection: deletionProtection,
	}

	idx, err := pc.CreateServerlessIndex(ctx, createRequest)
	if err != nil {
		msg.FailMsg("Failed to create serverless index %s: %s\n", style.Emphasis(options.name), err)
		exit.Error(err)
	}
	if options.json {
		json := text.IndentJSON(idx)
		pcio.Println(json)
		return
	}

	describeCommand := pcio.Sprintf("pc index describe --name %s", idx.Name)
	msg.SuccessMsg("Index %s created successfully. Run %s to check status. \n\n", style.Emphasis(idx.Name), style.Code(describeCommand))
	presenters.PrintDescribeIndexTable(idx)
}
