package index

import (
	"context"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/docslinks"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

// Abstracts the Pinecone Go SDK for unit testing (runCreateIndex)
type CreateIndexService interface {
	CreateServerlessIndex(ctx context.Context, req *pinecone.CreateServerlessIndexRequest) (*pinecone.Index, error)
	CreatePodIndex(ctx context.Context, req *pinecone.CreatePodIndexRequest) (*pinecone.Index, error)
	CreateIndexForModel(ctx context.Context, req *pinecone.CreateIndexForModelRequest) (*pinecone.Index, error)
	CreateBYOCIndex(ctx context.Context, req *pinecone.CreateBYOCIndexRequest) (*pinecone.Index, error)
}

type indexType string

const (
	indexTypeServerless indexType = "serverless"
	indexTypeIntegrated indexType = "integrated"
	indexTypeBYOC       indexType = "byoc"
	indexTypePod        indexType = "pod"
)

type createIndexOptions struct {
	// required for all index types
	name string

	// serverless only
	vectorType string

	// integrated only
	model           string
	fieldMap        map[string]string
	readParameters  map[string]string
	writeParameters map[string]string

	// BYOC only
	byocEnvironment string

	// pods only
	environment    string
	podType        string
	shards         int32
	replicas       int32
	metadataConfig []string

	// serverless & integrated
	cloud        string
	region       string
	readMode     string
	readNodeType string
	readShards   int32
	readReplicas int32

	// serverless & pods
	sourceCollection string

	// serverless & integrated & BYOC
	metadataSchema []string

	// optional for all index types
	dimension          int32
	metric             string
	deletionProtection string
	tags               map[string]string
	json               bool
}

var (
	createIndexHelp = help.LongF(`
		Create a new index with the specified configuration.
		
		You can specify the measure of similarity, the dimension of vectors to be stored, and which cloud
		provider to deploy with. You can also control whether the index is 'sparse' or 'dense',
		and any integrated embedding configuration you'd like to use.

		See: %s
	`, docslinks.DocsIndexCreate)

	createIndexExample = help.Examples(`
		# create a serverless index
		pc index create --name "my-index" --dimension 1536 --metric "cosine" --cloud "aws" --region "us-east-1"

		# create a pod index
		pc index create --name "my-index" --dimension 1536 --metric "cosine" --environment "us-east-1-aws" --pod-type "p1.x1" --shards 2 --replicas 2

		# create an integrated index
		pc index create --name "my-index" --dimension 1536 --metric "cosine" --cloud "aws" --region "us-east-1" --model "multilingual-e5-large" --field-map "text=chunk_text"
	`)
)

func NewCreateIndexCmd() *cobra.Command {
	options := createIndexOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new index with the specified configuration",
		Long:    createIndexHelp,
		Example: createIndexExample,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)

			idx, err := runCreateIndexCmd(ctx, cmd, pc, options)
			if err != nil {
				msg.FailMsg("Failed to create index: %s\n", err)
				exit.Error(err, "Failed to create index")
			}

			renderSuccessOutput(idx, options)
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "Name of index to create")
	_ = cmd.MarkFlagRequired("name")

	// Serverless & Pods
	cmd.Flags().StringVar(&options.sourceCollection, "source-collection", "", "When creating an index from a collection")

	// Serverless, BYOC, and Integrated
	cmd.Flags().StringSliceVar(&options.metadataSchema, "schema", []string{}, "Schema for the behavior of Pinecone's internal metadata index. By default, all metadata is indexed; when schema is present, only the fields provided will be indexed")

	// BYOC
	cmd.Flags().StringVar(&options.byocEnvironment, "byoc-environment", "", "BYOC environment to use for the index")

	// Serverless & Integrated
	cmd.Flags().StringVarP(&options.cloud, "cloud", "c", "", "Cloud provider where you would like to deploy your index")
	cmd.Flags().StringVarP(&options.region, "region", "r", "", "Cloud region where you would like to deploy your index")
	cmd.Flags().StringVar(&options.readMode, "read-mode", "", "The read capacity mode to use. One of: ondemand, dedicated. Defaults to ondemand. If configuring dedicated, you must also provide read-node-type, read-shards, and read-replicas")
	cmd.Flags().StringVar(&options.readNodeType, "read-node-type", "", "The type of machines to use. Available options: b1 and t1. t1 includes increased processing power and memory")
	cmd.Flags().Int32Var(&options.readShards, "read-shards", 0, "The number of shards to use. Shards determine the storage capacity of an index, with each shard providing 250 GB of storage")
	cmd.Flags().Int32Var(&options.readReplicas, "read-replicas", 0, "The number of replicas to use. Replicas duplicate the compute resources and data of an index, allowing higher query throughput and availability")

	// Serverless flags
	cmd.Flags().StringVarP(&options.vectorType, "vector-type", "v", "", "Vector type to use. One of: dense, sparse")

	// Pod flags
	cmd.Flags().StringVar(&options.environment, "environment", "", "Environment of the index to create")
	cmd.Flags().StringVar(&options.podType, "pod-type", "", "Type of pod to use")
	cmd.Flags().Int32Var(&options.shards, "shards", 1, "Shards of the index to create")
	cmd.Flags().Int32Var(&options.replicas, "replicas", 1, "Replicas of the index to create")
	cmd.Flags().StringSliceVar(&options.metadataConfig, "metadata-config", []string{}, "Metadata configuration to limit the fields that are indexed for search")

	// Integrated flags
	cmd.Flags().StringVar(&options.model, "model", "", "The name of the embedding model to use for the index")
	cmd.Flags().StringToStringVar(&options.fieldMap, "field-map", map[string]string{}, "Identifies the name of the text field from your document model that will be embedded")
	cmd.Flags().StringToStringVar(&options.readParameters, "read-parameters", map[string]string{}, "The read parameters for the embedding model")
	cmd.Flags().StringToStringVar(&options.writeParameters, "write-parameters", map[string]string{}, "The write parameters for the embedding model")

	// Optional flags - all index types
	cmd.Flags().Int32VarP(&options.dimension, "dimension", "d", 0, "Dimension of the index to create")
	cmd.Flags().StringVarP(&options.metric, "metric", "m", "cosine", "Metric to use. One of: cosine, euclidean, dotproduct")
	cmd.Flags().StringVar(&options.deletionProtection, "deletion-protection", "", "Whether to enable deletion protection for the index. One of: enabled, disabled")
	cmd.Flags().StringToStringVar(&options.tags, "tags", map[string]string{}, "Custom user tags to add to an index")

	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}

func runCreateIndexCmd(ctx context.Context, cmd *cobra.Command, service CreateIndexService, options createIndexOptions) (*pinecone.Index, error) {
	// validate and derive index type from arguments
	err := options.validate()
	if err != nil {
		return nil, err
	}
	idxType, err := options.deriveIndexType()
	if err != nil {
		return nil, err
	}

	// index tags
	var indexTags *pinecone.IndexTags
	if len(options.tags) > 0 {
		tags := pinecone.IndexTags(options.tags)
		indexTags = &tags
	}

	// read capacity configuration
	readCapacity, err := buildReadCapacityFromFlags(cmd, options.readMode, options.readNodeType, options.readShards, options.readReplicas)
	if err != nil {
		return nil, err
	}

	// created index
	var idx *pinecone.Index

	switch idxType {
	case indexTypeServerless:
		// create serverless index
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
			ReadCapacity:       readCapacity,
			Schema:             sdk.BuildMetadataSchema(options.metadataSchema),
		}

		idx, err = service.CreateServerlessIndex(ctx, &args)
		if err != nil {
			wrapped := pcio.Errorf("Failed to create serverless index %s: %w", style.Emphasis(options.name), err)
			return nil, wrapped
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

		idx, err = service.CreatePodIndex(ctx, &args)
		if err != nil {
			wrapped := pcio.Errorf("Failed to create pod index %s: %w", style.Emphasis(options.name), err)
			return nil, wrapped
		}
	case indexTypeIntegrated:
		// create integrated index
		readParams := toInterfaceMap(options.readParameters)
		writeParams := toInterfaceMap(options.writeParameters)

		readCapacity, err := buildReadCapacityFromFlags(cmd, options.readMode, options.readNodeType, options.readShards, options.readReplicas)
		if err != nil {
			return nil, err
		}

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
			Tags:         indexTags,
			ReadCapacity: readCapacity,
			Schema:       sdk.BuildMetadataSchema(options.metadataSchema),
		}

		idx, err = service.CreateIndexForModel(ctx, &args)
		if err != nil {
			wrapped := pcio.Errorf("Failed to create integrated index %s: %w", style.Emphasis(options.name), err)
			return nil, wrapped
		}
	case indexTypeBYOC:
		// create BYOC index
		args := pinecone.CreateBYOCIndexRequest{
			Name:               options.name,
			Environment:        options.byocEnvironment,
			Metric:             pointerOrNil(pinecone.IndexMetric(options.metric)),
			DeletionProtection: pointerOrNil(pinecone.DeletionProtection(options.deletionProtection)),
			Dimension:          pointerOrNil(options.dimension),
			Tags:               indexTags,
			Schema:             sdk.BuildMetadataSchema(options.metadataSchema),
		}

		idx, err = service.CreateBYOCIndex(ctx, &args)
		if err != nil {
			wrapped := pcio.Errorf("Failed to create BYOC index %s: %w", style.Emphasis(options.name), err)
			return nil, wrapped
		}
	default:
		err := pcio.Errorf("Error creating index: invalid index type")
		return nil, err
	}

	return idx, nil
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
	if c.byocEnvironment != "" {
		return indexTypeBYOC, nil
	}
	if c.environment != "" {
		return indexTypePod, nil
	}
	return "", pcio.Error("invalid index type. Please provide either environment, or cloud and region")
}

// Builds the ReadCapacityParams object based on the provided arguments
// "OnDemand" is the default with no explicit configuration. "Dedicated" requires nodeType, shards, and replicas
// for creating an index, or migrating to a dedicated index.
func buildReadCapacityFromFlags(cmd *cobra.Command, mode, nodeType string, shards, replicas int32) (*pinecone.ReadCapacityParams, error) {
	// only read flags that have been set by the user
	modeSet := cmd.Flags().Changed("read-mode")
	nodeSet := cmd.Flags().Changed("read-node-type")
	shardsSet := cmd.Flags().Changed("read-shards")
	replSet := cmd.Flags().Changed("read-replicas")

	var nodeTypePtr *string
	var shardsPtr *int32
	var replicasPtr *int32
	if nodeSet {
		nodeTypePtr = &nodeType
	}
	if shardsSet {
		shardsPtr = &shards
	}
	if replSet {
		replicasPtr = &replicas
	}

	// If no arguments are provided, pinecone.ReadCapacityParams should be nil
	if !modeSet && !nodeSet && !shardsSet && !replSet {
		return nil, nil
	}

	normMode := strings.ToLower(mode)
	// read-mode specifically requested
	if modeSet {
		switch normMode {
		case "ondemand":
			if nodeSet || shardsSet || replSet {
				return nil, pcio.Errorf("read-node-type, read-shards, and read-replicas are not supported with read-mode=ondemand")
			}
			return &pinecone.ReadCapacityParams{
				OnDemand: &pinecone.ReadCapacityOnDemandConfig{},
			}, nil
		case "dedicated":
			// continue
		default:
			return nil, pcio.Errorf("invalid read-mode")
		}
	} else { // read-mode not provided, return nil if no specific configuration values are passed
		if !nodeSet && !shardsSet && !replSet {
			return nil, nil
		}
	}

	// dedicated mode if ondemand mode was not requested
	return &pinecone.ReadCapacityParams{
		Dedicated: &pinecone.ReadCapacityDedicatedConfig{
			NodeType: nodeTypePtr,
			Scaling: &pinecone.ReadCapacityScaling{
				Manual: &pinecone.ReadCapacityManualScaling{
					Shards:   shardsPtr,
					Replicas: replicasPtr,
				},
			},
		},
	}, nil
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
