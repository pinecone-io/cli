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

type createPodOptions struct {
	name             string
	dimension        int32
	metric           string
	environment      string
	podType          string
	shards           int32
	replicas         int32
	sourceCollection string
	// metadataConfig   *PodSpecMetadataConfig

	json bool
}

func (o createPodOptions) isValidMetric() bool {
	switch pinecone.IndexMetric(o.metric) {
	case pinecone.Cosine, pinecone.Euclidean, pinecone.Dotproduct:
		return true
	default:
		return false
	}
}

func NewCreatePodCmd() *cobra.Command {
	options := createPodOptions{}

	cmd := &cobra.Command{
		Use:     "create-pod",
		Short:   "Create a pod index with the specified configuration",
		Example: "",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !options.isValidMetric() {
				return fmt.Errorf("metric must be one of [cosine, euclidean, dotproduct]")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runCreatePodCmd(cmd, options)
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to create")
	cmd.MarkFlagRequired("name")
	cmd.Flags().Int32VarP(&options.dimension, "dimension", "d", 0, "dimension of the index to create")
	cmd.MarkFlagRequired("dimension")
	cmd.Flags().StringVarP(&options.environment, "environment", "e", "", "environment of the index to create")
	cmd.MarkFlagRequired("environment")
	cmd.Flags().StringVarP(&options.podType, "pod_type", "t", "", "type of pod to use")
	cmd.MarkFlagRequired("pod_type")

	// Optional flags
	cmd.Flags().StringVarP(&options.metric, "metric", "m", "cosine", "metric to use. One of: cosine, euclidean, dotproduct")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().Int32VarP(&options.shards, "shards", "s", 1, "shards of the index to create")
	cmd.Flags().Int32VarP(&options.replicas, "replicas", "r", 1, "replicas of the index to create")
	cmd.Flags().StringVarP(&options.sourceCollection, "source_collection", "c", "", "When creating a pod index using data from a collection, the name of the source collection")
	cmd.MarkFlagRequired("sourceCollection")

	return cmd
}

func runCreatePodCmd(cmd *cobra.Command, options createPodOptions) {
	ctx := context.Background()
	pc := client.NewPineconeClient()

	createRequest := &pinecone.CreatePodIndexRequest{
		Name:        options.name,
		Metric:      pinecone.IndexMetric(options.metric),
		Dimension:   options.dimension,
		Environment: options.environment,
		PodType:     options.podType,
		Shards:      options.shards,
		Replicas:    options.replicas,
	}

	idx, err := pc.CreatePodIndex(ctx, createRequest)
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
