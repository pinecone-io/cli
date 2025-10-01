package index

import (
	"context"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type createPodOptions struct {
	name               string
	dimension          int32
	metric             string
	environment        string
	podType            string
	shards             int32
	replicas           int32
	sourceCollection   string
	deletionProtection string
	// metadataConfig   *PodSpecMetadataConfig

	json bool
}

func NewCreatePodCmd() *cobra.Command {
	options := createPodOptions{}

	cmd := &cobra.Command{
		Use:   "create-pod",
		Short: "Create a pod index with the specified configuration",
		Example: help.Examples(`
			pc index create-pod --name "my-index" --dimension 1536 --metric "cosine" --environment "us-east-1-aws" --pod-type "p1.x1" --shards 2 --replicas 2
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runCreatePodCmd(options)
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to create")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().Int32VarP(&options.dimension, "dimension", "d", 0, "dimension of the index to create")
	_ = cmd.MarkFlagRequired("dimension")
	cmd.Flags().StringVarP(&options.environment, "environment", "e", "", "environment of the index to create")
	_ = cmd.MarkFlagRequired("environment")
	cmd.Flags().StringVarP(&options.podType, "pod_type", "t", "", "type of pod to use")
	_ = cmd.MarkFlagRequired("pod_type")

	// Optional flags
	cmd.Flags().StringVarP(&options.metric, "metric", "m", "cosine", "metric to use. One of: cosine, euclidean, dotproduct")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().Int32VarP(&options.shards, "shards", "s", 1, "shards of the index to create")
	cmd.Flags().Int32VarP(&options.replicas, "replicas", "r", 1, "replicas of the index to create")
	cmd.Flags().StringVarP(&options.sourceCollection, "source_collection", "c", "", "When creating a pod index using data from a collection, the name of the source collection")
	cmd.Flags().StringVarP(&options.deletionProtection, "deletion_protection", "p", "", "Whether to enable deletion protection for the index")
	_ = cmd.MarkFlagRequired("sourceCollection")

	return cmd
}

func runCreatePodCmd(options createPodOptions) {
	ctx := context.Background()
	pc := sdk.NewPineconeClient()

	// Deprecation warning
	pcio.Fprintf(os.Stderr, "⚠️  Warning: The '%s' command is deprecated. Please use '%s' instead.", style.Code("index create-pod"), style.Code("index create"))

	metric := pinecone.IndexMetric(options.metric)
	deletionProtection := pinecone.DeletionProtection(options.deletionProtection)
	createRequest := &pinecone.CreatePodIndexRequest{
		Name:               options.name,
		Metric:             &metric,
		Dimension:          options.dimension,
		Environment:        options.environment,
		PodType:            options.podType,
		Shards:             options.shards,
		Replicas:           options.replicas,
		DeletionProtection: &deletionProtection,
	}

	idx, err := pc.CreatePodIndex(ctx, createRequest)
	if err != nil {
		msg.FailMsg("Failed to create index %s: %s\n", style.Emphasis(options.name), err)
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
