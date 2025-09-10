package index

import (
	"context"
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/docslinks"
	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/index"
	indexpresenters "github.com/pinecone-io/cli/internal/pkg/utils/index/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/interactive"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type createIndexOptions struct {
	CreateOptions index.CreateOptions
	json          bool
}

func NewCreateIndexCmd() *cobra.Command {
	options := createIndexOptions{}

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new index with the specified configuration",
		Long: heredoc.Docf(`
		The %s command creates a new index with the specified configuration. There are several different types of indexes
		you can create depending on the configuration provided:

			- Serverless (dense or sparse)
			- Integrated 
			- Pod

		For detailed documentation, see:
		%s
		`, style.Code("pc index create"), style.URL(docslinks.DocsIndexCreate)),
		Example: heredoc.Doc(`
		# create a serverless index
		$ pc index create my-index --dimension 1536 --metric cosine --cloud aws --region us-east-1

		# create a pod index
		$ pc index create my-index --dimension 1536 --metric cosine --environment us-east-1-aws --pod-type p1.x1 --shards 2 --replicas 2

		# create an integrated index
		$ pc index create my-index --dimension 1536 --metric cosine --cloud aws --region us-east-1 --model multilingual-e5-large --field_map text=chunk_text
		`),
		Args:         index.ValidateIndexNameArgs,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			options.CreateOptions.Name.Value = args[0]
			runCreateIndexCmd(options, cmd, args)
		},
	}

	// index type flags
	cmd.Flags().BoolVar(&options.CreateOptions.Serverless.Value, "serverless", false, "Create a serverless index (default)")
	cmd.Flags().BoolVar(&options.CreateOptions.Pod.Value, "pod", false, "Create a pod index")

	// Serverless & Pods
	cmd.Flags().StringVar(&options.CreateOptions.SourceCollection.Value, "source_collection", "", "When creating an index from a collection")

	// Serverless & Integrated
	cmd.Flags().StringVarP(&options.CreateOptions.Cloud.Value, "cloud", "c", "", "Cloud provider where you would like to deploy your index")
	cmd.Flags().StringVarP(&options.CreateOptions.Region.Value, "region", "r", "", "Cloud region where you would like to deploy your index")

	// Serverless flags
	cmd.Flags().StringVarP(&options.CreateOptions.VectorType.Value, "vector_type", "v", "", "Vector type to use. One of: dense, sparse")

	// Pod flags
	cmd.Flags().StringVar(&options.CreateOptions.Environment.Value, "environment", "", "Environment of the index to create")
	cmd.Flags().StringVar(&options.CreateOptions.PodType.Value, "pod_type", "", "Type of pod to use")
	cmd.Flags().Int32Var(&options.CreateOptions.Shards.Value, "shards", 1, "Shards of the index to create")
	cmd.Flags().Int32Var(&options.CreateOptions.Replicas.Value, "replicas", 1, "Replicas of the index to create")
	cmd.Flags().StringSliceVar(&options.CreateOptions.MetadataConfig.Value, "metadata_config", []string{}, "Metadata configuration to limit the fields that are indexed for search")

	// Integrated flags
	cmd.Flags().StringVar(&options.CreateOptions.Model.Value, "model", "", "The name of the embedding model to use for the index")
	cmd.Flags().StringToStringVar(&options.CreateOptions.FieldMap.Value, "field_map", map[string]string{}, "Identifies the name of the text field from your document model that will be embedded")
	cmd.Flags().StringToStringVar(&options.CreateOptions.ReadParameters.Value, "read_parameters", map[string]string{}, "The read parameters for the embedding model")
	cmd.Flags().StringToStringVar(&options.CreateOptions.WriteParameters.Value, "write_parameters", map[string]string{}, "The write parameters for the embedding model")

	// Optional flags
	cmd.Flags().Int32VarP(&options.CreateOptions.Dimension.Value, "dimension", "d", 0, "Dimension of the index to create")
	cmd.Flags().StringVarP(&options.CreateOptions.Metric.Value, "metric", "m", "", "Metric to use. One of: cosine, euclidean, dotproduct")
	cmd.Flags().StringVar(&options.CreateOptions.DeletionProtection.Value, "deletion_protection", "", "Whether to enable deletion protection for the index. One of: enabled, disabled")
	cmd.Flags().StringToStringVar(&options.CreateOptions.Tags.Value, "tags", map[string]string{}, "Custom user tags to add to an index")

	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}

func runCreateIndexCmd(options createIndexOptions, cmd *cobra.Command, args []string) {

	// validationErrors := index.ValidateCreateOptions(options.CreateOptions)
	// if len(validationErrors) > 0 {
	// 	msg.FailMsgMultiLine(validationErrors...)
	// 	exit.Error(errors.New(validationErrors[0])) // Use first error for exit code
	// }

	inferredOptions := index.InferredCreateOptions(options.CreateOptions)
	validationErrors := index.ValidateCreateOptions(inferredOptions)
	if len(validationErrors) > 0 {
		msg.FailMsgMultiLine(validationErrors...)
		exit.Error(errors.New(validationErrors[0])) // Use first error for exit code
	}

	// Print preview of what will be created
	pcio.Println()
	pcio.Printf("%s\n\n",
		pcio.Sprintf("Creating %s index %s with the following configuration:",
			style.Emphasis(string(inferredOptions.GetSpec())),
			style.ResourceName(inferredOptions.Name.Value),
		),
	)

	indexpresenters.PrintIndexCreateConfigTable(&inferredOptions)

	// Ask for user confirmation
	question := "Is this configuration correct? Do you want to proceed with creating the index?"
	if !interactive.GetConfirmation(question) {
		pcio.Println(style.InfoMsg("Index creation cancelled."))
		return
	}

	// index tags
	var indexTags *pinecone.IndexTags
	if len(inferredOptions.Tags.Value) > 0 {
		tags := pinecone.IndexTags(inferredOptions.Tags.Value)
		indexTags = &tags
	}

	// created index
	var idx *pinecone.Index
	var err error
	ctx := context.Background()
	pc := sdk.NewPineconeClient()

	switch inferredOptions.GetCreateFlow() {
	case index.Serverless:
		// create serverless index
		req := pinecone.CreateServerlessIndexRequest{
			Name:               inferredOptions.Name.Value,
			Cloud:              pinecone.Cloud(inferredOptions.Cloud.Value),
			Region:             inferredOptions.Region.Value,
			Metric:             pointerOrNil(pinecone.IndexMetric(inferredOptions.Metric.Value)),
			DeletionProtection: pointerOrNil(pinecone.DeletionProtection(inferredOptions.DeletionProtection.Value)),
			Dimension:          pointerOrNil(inferredOptions.Dimension.Value),
			VectorType:         pointerOrNil(inferredOptions.VectorType.Value),
			Tags:               indexTags,
			SourceCollection:   pointerOrNil(inferredOptions.SourceCollection.Value),
		}

		idx, err = pc.CreateServerlessIndex(ctx, &req)
		if err != nil {
			errorutil.HandleIndexAPIError(err, cmd, args)
			exit.Error(err)
		}
	case index.Pod:
		// create pod index
		var metadataConfig *pinecone.PodSpecMetadataConfig
		if len(inferredOptions.MetadataConfig.Value) > 0 {
			metadataConfig = &pinecone.PodSpecMetadataConfig{
				Indexed: &inferredOptions.MetadataConfig.Value,
			}
		}
		req := pinecone.CreatePodIndexRequest{
			Name:               inferredOptions.Name.Value,
			Dimension:          inferredOptions.Dimension.Value,
			Environment:        inferredOptions.Environment.Value,
			PodType:            inferredOptions.PodType.Value,
			Shards:             inferredOptions.Shards.Value,
			Replicas:           inferredOptions.Replicas.Value,
			Metric:             pointerOrNil(pinecone.IndexMetric(inferredOptions.Metric.Value)),
			DeletionProtection: pointerOrNil(pinecone.DeletionProtection(inferredOptions.DeletionProtection.Value)),
			SourceCollection:   pointerOrNil(inferredOptions.SourceCollection.Value),
			Tags:               indexTags,
			MetadataConfig:     metadataConfig,
		}

		idx, err = pc.CreatePodIndex(ctx, &req)
		if err != nil {
			errorutil.HandleIndexAPIError(err, cmd, args)
			exit.Error(err)
		}
	case index.Integrated:
		// create integrated index
		readParams := toInterfaceMap(inferredOptions.ReadParameters.Value)
		writeParams := toInterfaceMap(inferredOptions.WriteParameters.Value)

		req := pinecone.CreateIndexForModelRequest{
			Name:               inferredOptions.Name.Value,
			Cloud:              pinecone.Cloud(inferredOptions.Cloud.Value),
			Region:             inferredOptions.Region.Value,
			DeletionProtection: pointerOrNil(pinecone.DeletionProtection(inferredOptions.DeletionProtection.Value)),
			Embed: pinecone.CreateIndexForModelEmbed{
				Model:           inferredOptions.Model.Value,
				FieldMap:        toInterfaceMap(inferredOptions.FieldMap.Value),
				ReadParameters:  &readParams,
				WriteParameters: &writeParams,
			},
		}

		idx, err = pc.CreateIndexForModel(ctx, &req)
		if err != nil {
			errorutil.HandleIndexAPIError(err, cmd, args)
			exit.Error(err)
		}
	default:
		err := pcio.Errorf("invalid index type")
		log.Error().Err(err).Msg("Error creating index")
		exit.Error(err)
	}

	renderSuccessOutput(idx, options.json)
}

func renderSuccessOutput(idx *pinecone.Index, jsonOutput bool) {
	if jsonOutput {
		json := text.IndentJSON(idx)
		pcio.Println(json)
		return
	}

	msg.SuccessMsg("Index %s created successfully.", style.ResourceName(idx.Name))

	indexpresenters.PrintDescribeIndexTable(idx)

	describeCommand := pcio.Sprintf("pc index describe %s", idx.Name)
	hint := fmt.Sprintf("Run %s at any time to check the status. \n\n", style.Code(describeCommand))
	pcio.Println(style.Hint(hint))
}

func pointerOrNil[T comparable](value T) *T {
	var zero T // set to zero-value of generic type T
	if value == zero {
		return nil
	}
	return &value
}

func toInterfaceMap(in map[string]string) map[string]any {
	if in == nil {
		return nil
	}

	interfaceMap := make(map[string]any, len(in))
	for k, v := range in {
		interfaceMap[k] = v
	}
	return interfaceMap
}
