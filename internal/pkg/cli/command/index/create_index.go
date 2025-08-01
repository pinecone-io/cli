package index

import (
	"context"
	"errors"
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

type indexType string

const (
	indexTypeServerless indexType = "serverless"
	indexTypeIntegrated indexType = "integrated"
	indexTypePod        indexType = "pod"
)

type createIndexOptions struct {
	// required for all index types
	name string

	// serverless only
	vectorType string

	// serverless & integrated
	cloud  string
	region string

	// serverless & pods
	sourceCollection string

	// pods only
	environment    string
	podType        string
	shards         int32
	replicas       int32
	metadataConfig []string

	// integrated only
	model           string
	fieldMap        map[string]string
	readParameters  map[string]string
	writeParameters map[string]string

	// optional for all index types
	dimension          int32
	metric             string
	deletionProtection string
	tags               map[string]string

	json bool
}

func NewCreateIndexCmd() *cobra.Command {
	options := createIndexOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new index with the specified configuration",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			runCreateIndexCmd(options)
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to create")
	cmd.MarkFlagRequired("name")

	// Serverless & Pods
	cmd.Flags().StringVar(&options.sourceCollection, "source_collection", "", "when creating an index from a collection")

	// Serverless & Integrated
	cmd.Flags().StringVarP(&options.cloud, "cloud", "c", "", "cloud provider where you would like to deploy your index")
	cmd.Flags().StringVarP(&options.region, "region", "r", "", "cloud region where you would like to deploy your index")

	// Serverless flags
	cmd.Flags().StringVarP(&options.vectorType, "vector_type", "v", "", "vector type to use. One of: dense, sparse")

	// Pod flags
	cmd.Flags().StringVar(&options.environment, "environment", "", "environment of the index to create")
	cmd.Flags().StringVar(&options.podType, "pod_type", "", "type of pod to use")
	cmd.Flags().Int32Var(&options.shards, "shards", 1, "shards of the index to create")
	cmd.Flags().Int32Var(&options.replicas, "replicas", 1, "replicas of the index to create")
	cmd.Flags().StringSliceVar(&options.metadataConfig, "metadata_config", []string{}, "metadata configuration to limit the fields that are indexed for search")

	// Integrated flags
	cmd.Flags().StringVar(&options.model, "model", "", "the name of the embedding model to use for the index")
	cmd.Flags().StringToStringVar(&options.fieldMap, "field_map", map[string]string{}, "identifies the name of the text field from your document model that will be embedded")
	cmd.Flags().StringToStringVar(&options.readParameters, "read_parameters", map[string]string{}, "the read parameters for the embedding model")
	cmd.Flags().StringToStringVar(&options.writeParameters, "write_parameters", map[string]string{}, "the write parameters for the embedding model")

	// Optional flags
	cmd.Flags().Int32VarP(&options.dimension, "dimension", "d", 0, "dimension of the index to create")
	cmd.Flags().StringVarP(&options.metric, "metric", "m", "cosine", "metric to use. One of: cosine, euclidean, dotproduct")
	cmd.Flags().StringVar(&options.deletionProtection, "deletion_protection", "", "whether to enable deletion protection for the index. One of: enabled, disabled")
	cmd.Flags().StringToStringVar(&options.tags, "tags", map[string]string{}, "custom user tags to add to an index")

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func runCreateIndexCmd(options createIndexOptions) {
	ctx := context.Background()
	pc := sdk.NewPineconeClient()

	// validate and derive index type from arguments
	err := options.validate()
	if err != nil {
		exit.Error(err)
		return
	}
	idxType, err := options.deriveIndexType()
	if err != nil {
		exit.Error(err)
		return
	}

	// index tags
	var indexTags *pinecone.IndexTags
	if len(options.tags) > 0 {
		tags := pinecone.IndexTags(options.tags)
		indexTags = &tags
	}

	// created index
	var idx *pinecone.Index

	switch idxType {
	case indexTypeServerless:
		// create serverless index
		fmt.Printf("dimension %d\n", options.dimension)
		args := pinecone.CreateServerlessIndexRequest{
			Name:               options.name,
			Cloud:              pinecone.Cloud(options.cloud),
			Region:             options.region,
			Metric:             pointerOrNil(pinecone.IndexMetric(options.metric)),
			DeletionProtection: pointerOrNil(pinecone.DeletionProtection(options.deletionProtection)),
			Dimension:          pointerOrNil(options.dimension),
			VectorType:         pointerOrNil(options.vectorType),
			Tags:               indexTags,
			SourceCollection:   pointerOrNil(options.sourceCollection),
		}

		idx, err = pc.CreateServerlessIndex(ctx, &args)
		if err != nil {
			msg.FailMsg("Failed to create serverless index %s: %s\n", style.Emphasis(options.name), err)
			exit.Error(err)
		}
	case indexTypePod:
		// create pod index
		var metadataConfig *pinecone.PodSpecMetadataConfig
		if len(options.metadataConfig) > 0 {
			metadataConfig = &pinecone.PodSpecMetadataConfig{
				Indexed: &options.metadataConfig,
			}
		}
		args := pinecone.CreatePodIndexRequest{
			Name:               options.name,
			Dimension:          options.dimension,
			Environment:        options.environment,
			PodType:            options.podType,
			Shards:             options.shards,
			Replicas:           options.replicas,
			Metric:             pointerOrNil(pinecone.IndexMetric(options.metric)),
			DeletionProtection: pointerOrNil(pinecone.DeletionProtection(options.deletionProtection)),
			SourceCollection:   pointerOrNil(options.sourceCollection),
			Tags:               indexTags,
			MetadataConfig:     metadataConfig,
		}

		idx, err = pc.CreatePodIndex(ctx, &args)
		if err != nil {
			msg.FailMsg("Failed to create pod index %s: %s\n", style.Emphasis(options.name), err)
			exit.Error(err)
		}
	case indexTypeIntegrated:
		// create integrated index
		readParams := toInterfaceMap(options.readParameters)
		writeParams := toInterfaceMap(options.writeParameters)

		args := pinecone.CreateIndexForModelRequest{
			Name:               options.name,
			Cloud:              pinecone.Cloud(options.cloud),
			Region:             options.region,
			DeletionProtection: pointerOrNil(pinecone.DeletionProtection(options.deletionProtection)),
			Embed: pinecone.CreateIndexForModelEmbed{
				Model:           options.model,
				FieldMap:        toInterfaceMap(options.fieldMap),
				ReadParameters:  &readParams,
				WriteParameters: &writeParams,
			},
		}

		idx, err = pc.CreateIndexForModel(ctx, &args)
		if err != nil {
			msg.FailMsg("Failed to create integrated index %s: %s\n", style.Emphasis(options.name), err)
			exit.Error(err)
		}
	default:
		err := pcio.Errorf("invalid index type")
		log.Error().Err(err).Msg("Error creating index")
		exit.Error(err)
	}

	renderSuccessOutput(idx, options)
}

func renderSuccessOutput(idx *pinecone.Index, options createIndexOptions) {
	if options.json {
		json := text.IndentJSON(idx)
		pcio.Println(json)
		return
	}

	describeCommand := pcio.Sprintf("pc index describe --name %s", idx.Name)
	msg.SuccessMsg("Index %s created successfully. Run %s to check status. \n\n", style.Emphasis(idx.Name), style.Code(describeCommand))
	presenters.PrintDescribeIndexTable(idx)
}

// validate specific input params
func (c *createIndexOptions) validate() error {
	// name required for all index types
	if c.name == "" {
		err := pcio.Errorf("name is required")
		log.Error().Err(err).Msg("Error creating index")
		return err
	}

	// environment and cloud/region cannot be provided together
	if c.cloud != "" && c.region != "" && c.environment != "" {
		err := pcio.Errorf("cloud, region, and environment cannot be provided together")
		log.Error().Err(err).Msg("Error creating index")
		return err
	}

	return nil
}

// determine the type of index being created based on high level input params
func (c *createIndexOptions) deriveIndexType() (indexType, error) {
	if c.cloud != "" && c.region != "" {
		if c.model != "" {
			return indexTypeIntegrated, nil
		} else {
			return indexTypeServerless, nil
		}
	}
	if c.environment != "" {
		return indexTypePod, nil
	}
	return "", errors.New("invalid index type. Please provide either environment, or cloud and region")
}

func pointerOrNil[T comparable](value T) *T {
	var zero T // set to zero-value of generic type T
	if value == zero {
		return nil
	}
	return &value
}

func toInterfaceMap(in map[string]string) map[string]interface{} {
	if in == nil {
		return nil
	}

	interfaceMap := make(map[string]interface{}, len(in))
	for k, v := range in {
		interfaceMap[k] = v
	}
	return interfaceMap
}
