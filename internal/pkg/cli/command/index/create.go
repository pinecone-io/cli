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
			options.CreateOptions.Name = args[0]
			runCreateIndexCmd(options, cmd, args)
		},
	}

	// index type flags
	cmd.Flags().BoolVar(&options.CreateOptions.Serverless, "serverless", false, "Create a serverless index (default)")
	cmd.Flags().BoolVar(&options.CreateOptions.Pod, "pod", false, "Create a pod index")

	// Serverless & Pods
	cmd.Flags().StringVar(&options.CreateOptions.SourceCollection, "source_collection", "", "When creating an index from a collection")

	// Serverless & Integrated
	cmd.Flags().StringVarP(&options.CreateOptions.Cloud, "cloud", "c", "", "Cloud provider where you would like to deploy your index")
	cmd.Flags().StringVarP(&options.CreateOptions.Region, "region", "r", "", "Cloud region where you would like to deploy your index")

	// Serverless flags
	cmd.Flags().StringVarP(&options.CreateOptions.VectorType, "vector_type", "v", "", "Vector type to use. One of: dense, sparse")

	// Pod flags
	cmd.Flags().StringVar(&options.CreateOptions.Environment, "environment", "", "Environment of the index to create")
	cmd.Flags().StringVar(&options.CreateOptions.PodType, "pod_type", "", "Type of pod to use")
	cmd.Flags().Int32Var(&options.CreateOptions.Shards, "shards", 1, "Shards of the index to create")
	cmd.Flags().Int32Var(&options.CreateOptions.Replicas, "replicas", 1, "Replicas of the index to create")
	cmd.Flags().StringSliceVar(&options.CreateOptions.MetadataConfig, "metadata_config", []string{}, "Metadata configuration to limit the fields that are indexed for search")

	// Integrated flags
	cmd.Flags().StringVar(&options.CreateOptions.Model, "model", "", "The name of the embedding model to use for the index")
	cmd.Flags().StringToStringVar(&options.CreateOptions.FieldMap, "field_map", map[string]string{}, "Identifies the name of the text field from your document model that will be embedded")
	cmd.Flags().StringToStringVar(&options.CreateOptions.ReadParameters, "read_parameters", map[string]string{}, "The read parameters for the embedding model")
	cmd.Flags().StringToStringVar(&options.CreateOptions.WriteParameters, "write_parameters", map[string]string{}, "The write parameters for the embedding model")

	// Optional flags
	cmd.Flags().Int32VarP(&options.CreateOptions.Dimension, "dimension", "d", 0, "Dimension of the index to create")
	cmd.Flags().StringVarP(&options.CreateOptions.Metric, "metric", "m", "", "Metric to use. One of: cosine, euclidean, dotproduct")
	cmd.Flags().StringVar(&options.CreateOptions.DeletionProtection, "deletion_protection", "", "Whether to enable deletion protection for the index. One of: enabled, disabled")
	cmd.Flags().StringToStringVar(&options.CreateOptions.Tags, "tags", map[string]string{}, "Custom user tags to add to an index")

	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}

func runCreateIndexCmd(options createIndexOptions, cmd *cobra.Command, args []string) {

	validationErrors := index.ValidateCreateOptions(options.CreateOptions)
	if len(validationErrors) > 0 {
		msg.FailMsgMultiLine(validationErrors...)
		exit.Error(errors.New(validationErrors[0])) // Use first error for exit code
	}

	// Print preview of what will be created
	pcio.Println()
	pcio.Printf("%s\n\n",
		pcio.Sprintf("Creating %s index %s with the following configuration:",
			style.Emphasis(string(options.CreateOptions.GetSpec())),
			style.ResourceName(options.CreateOptions.Name),
		),
	)

	indexpresenters.PrintIndexCreateConfigTable(&options.CreateOptions)

	// Ask for user confirmation
	question := "Is this configuration correct? Do you want to proceed with creating the index?"
	if !interactive.GetConfirmation(question) {
		pcio.Println(style.InfoMsg("Index creation cancelled."))
		return
	}

	// index tags
	var indexTags *pinecone.IndexTags
	if len(options.CreateOptions.Tags) > 0 {
		tags := pinecone.IndexTags(options.CreateOptions.Tags)
		indexTags = &tags
	}

	// created index
	var idx *pinecone.Index
	var err error
	ctx := context.Background()
	pc := sdk.NewPineconeClient()

	switch options.CreateOptions.GetSpec() {
	case index.IndexSpecServerless:
		// create serverless index
		req := pinecone.CreateServerlessIndexRequest{
			Name:               options.CreateOptions.Name,
			Cloud:              pinecone.Cloud(options.CreateOptions.Cloud),
			Region:             options.CreateOptions.Region,
			Metric:             pointerOrNil(pinecone.IndexMetric(options.CreateOptions.Metric)),
			DeletionProtection: pointerOrNil(pinecone.DeletionProtection(options.CreateOptions.DeletionProtection)),
			Dimension:          pointerOrNil(options.CreateOptions.Dimension),
			VectorType:         pointerOrNil(options.CreateOptions.VectorType),
			Tags:               indexTags,
			SourceCollection:   pointerOrNil(options.CreateOptions.SourceCollection),
		}

		idx, err = pc.CreateServerlessIndex(ctx, &req)
		if err != nil {
			errorutil.HandleIndexAPIError(err, cmd, args)
			exit.Error(err)
		}
	case index.IndexSpecPod:
		// create pod index
		var metadataConfig *pinecone.PodSpecMetadataConfig
		if len(options.CreateOptions.MetadataConfig) > 0 {
			metadataConfig = &pinecone.PodSpecMetadataConfig{
				Indexed: &options.CreateOptions.MetadataConfig,
			}
		}
		req := pinecone.CreatePodIndexRequest{
			Name:               options.CreateOptions.Name,
			Dimension:          options.CreateOptions.Dimension,
			Environment:        options.CreateOptions.Environment,
			PodType:            options.CreateOptions.PodType,
			Shards:             options.CreateOptions.Shards,
			Replicas:           options.CreateOptions.Replicas,
			Metric:             pointerOrNil(pinecone.IndexMetric(options.CreateOptions.Metric)),
			DeletionProtection: pointerOrNil(pinecone.DeletionProtection(options.CreateOptions.DeletionProtection)),
			SourceCollection:   pointerOrNil(options.CreateOptions.SourceCollection),
			Tags:               indexTags,
			MetadataConfig:     metadataConfig,
		}

		idx, err = pc.CreatePodIndex(ctx, &req)
		if err != nil {
			errorutil.HandleIndexAPIError(err, cmd, args)
			exit.Error(err)
		}
	// case indexTypeIntegrated:
	// 	// create integrated index
	// 	readParams := toInterfaceMap(options.readParameters)
	// 	writeParams := toInterfaceMap(options.writeParameters)

	// 	req := pinecone.CreateIndexForModelRequest{
	// 		Name:               options.name,
	// 		Cloud:              pinecone.Cloud(options.cloud),
	// 		Region:             options.region,
	// 		DeletionProtection: pointerOrNil(pinecone.DeletionProtection(options.deletionProtection)),
	// 		Embed: pinecone.CreateIndexForModelEmbed{
	// 			Model:           options.model,
	// 			FieldMap:        toInterfaceMap(options.fieldMap),
	// 			ReadParameters:  &readParams,
	// 			WriteParameters: &writeParams,
	// 		},
	// 	}

	// 	idx, err = pc.CreateIndexForModel(ctx, &req)
	// 	if err != nil {
	// 		errorutil.HandleIndexAPIError(err, cmd, args)
	// 		exit.Error(err)
	// 	}
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
