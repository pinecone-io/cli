package index

import (
	"context"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/docslinks"
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
	// index type flags
	serverless bool
	pod        bool
	integrated bool

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
		Use:   "create [name]",
		Short: "Create a new index with the specified configuration",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("index name is required")
			}
			return nil
		},
		Long: heredoc.Docf(`
		The %s command creates a new index with the specified configuration. By default, it creates a serverless index
		with sensible defaults. You can explicitly specify the index type using the appropriate flag:

			- Serverless (default): Use --serverless or no flag
			  %s
			- Pod: Use --pod flag
			  %s
			- Integrated: Use --integrated flag
			  %s

		`, style.Code("pc index create"),
			style.URL(docslinks.DocsIndexCreate),
			style.URL(docslinks.DocsPodTypes),
			style.URL(docslinks.DocsIntegratedEmbedding)),
		Example: heredoc.Doc(`
		# create a serverless index (default)
		$ pc index create my-index

		# create a serverless index with custom configuration
		$ pc index create my-index --dimension 768 --metric euclidean --cloud aws --region us-east-1

		# create a pod index
		$ pc index create my-index --pod --environment us-east-1-aws --pod_type p1.x1 --shards 2 --replicas 2

		# create an integrated index
		$ pc index create my-index --integrated --cloud aws --region us-east-1 --model multilingual-e5-large --field_map text=chunk_text
		`),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			runCreateIndexCmd(name, options)
		},
	}

	// Index type flags
	cmd.Flags().BoolVar(&options.serverless, "serverless", false, "Create a serverless index (default)")
	cmd.Flags().BoolVar(&options.pod, "pod", false, "Create a pod index")
	cmd.Flags().BoolVar(&options.integrated, "integrated", false, "Create an integrated index")

	// Serverless & Pods
	cmd.Flags().StringVar(&options.sourceCollection, "source_collection", "", "When creating an index from a collection")

	// Serverless & Integrated
	cmd.Flags().StringVarP(&options.cloud, "cloud", "c", "", "Cloud provider where you would like to deploy your index")
	cmd.Flags().StringVarP(&options.region, "region", "r", "", "Cloud region where you would like to deploy your index")

	// Serverless flags
	cmd.Flags().StringVarP(&options.vectorType, "vector_type", "v", "", "Vector type to use. One of: dense, sparse")

	// Pod flags
	cmd.Flags().StringVar(&options.environment, "environment", "", "Environment of the index to create")
	cmd.Flags().StringVar(&options.podType, "pod_type", "", "Type of pod to use")
	cmd.Flags().Int32Var(&options.shards, "shards", 1, "Shards of the index to create")
	cmd.Flags().Int32Var(&options.replicas, "replicas", 1, "Replicas of the index to create")
	cmd.Flags().StringSliceVar(&options.metadataConfig, "metadata_config", []string{}, "Metadata configuration to limit the fields that are indexed for search")

	// Integrated flags
	cmd.Flags().StringVar(&options.model, "model", "", "The name of the embedding model to use for the index")
	cmd.Flags().StringToStringVar(&options.fieldMap, "field_map", map[string]string{}, "Identifies the name of the text field from your document model that will be embedded")
	cmd.Flags().StringToStringVar(&options.readParameters, "read_parameters", map[string]string{}, "The read parameters for the embedding model")
	cmd.Flags().StringToStringVar(&options.writeParameters, "write_parameters", map[string]string{}, "The write parameters for the embedding model")

	// Optional flags
	cmd.Flags().Int32VarP(&options.dimension, "dimension", "d", 1536, "Dimension of the index to create")
	cmd.Flags().StringVarP(&options.metric, "metric", "m", "cosine", "Metric to use. One of: cosine, euclidean, dotproduct")
	cmd.Flags().StringVar(&options.deletionProtection, "deletion_protection", "", "Whether to enable deletion protection for the index. One of: enabled, disabled")
	cmd.Flags().StringToStringVar(&options.tags, "tags", map[string]string{}, "Custom user tags to add to an index")

	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}

func runCreateIndexCmd(name string, options createIndexOptions) {
	ctx := context.Background()
	pc := sdk.NewPineconeClient()

	// validate and derive index type from arguments
	err := options.validate()
	if err != nil {
		msg.FailMsg("Validation failed: %s", err)
		exit.Error(err)
		return
	}
	idxType, err := options.deriveIndexType()
	if err != nil {
		msg.FailMsg("Configuration error: %s", err)
		exit.Error(err)
		return
	}

	// Print preview of what will be created
	printCreatePreview(name, options, idxType)

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
		args := pinecone.CreateServerlessIndexRequest{
			Name:               name,
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
			msg.FailMsg("Failed to create serverless index %s: %s\n", style.Emphasis(name), err)
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
			Name:               name,
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
			msg.FailMsg("Failed to create pod index %s: %s\n", style.Emphasis(name), err)
			exit.Error(err)
		}
	case indexTypeIntegrated:
		// create integrated index
		readParams := toInterfaceMap(options.readParameters)
		writeParams := toInterfaceMap(options.writeParameters)

		args := pinecone.CreateIndexForModelRequest{
			Name:               name,
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
			msg.FailMsg("Failed to create integrated index %s: %s\n", style.Emphasis(name), err)
			exit.Error(err)
		}
	default:
		err := pcio.Errorf("invalid index type")
		log.Error().Err(err).Msg("Error creating index")
		exit.Error(err)
	}

	renderSuccessOutput(idx, options)
}

// printCreatePreview prints a preview of the index configuration that will be created
func printCreatePreview(name string, options createIndexOptions, idxType indexType) {
	pcio.Println()
	pcio.Printf("Creating %s index '%s' with the following configuration:\n\n", style.Emphasis(string(idxType)), style.Emphasis(name))

	writer := presenters.NewTabWriter()
	log.Debug().Str("name", name).Msg("Printing index creation preview")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Name\t%s\n", name)
	pcio.Fprintf(writer, "Type\t%s\n", string(idxType))

	if options.dimension > 0 {
		pcio.Fprintf(writer, "Dimension\t%d\n", options.dimension)
	}

	pcio.Fprintf(writer, "Metric\t%s\n", options.metric)

	if options.deletionProtection != "" {
		pcio.Fprintf(writer, "Deletion Protection\t%s\n", options.deletionProtection)
	}

	if options.vectorType != "" {
		pcio.Fprintf(writer, "Vector Type\t%s\n", options.vectorType)
	}

	pcio.Fprintf(writer, "\t\n")

	switch idxType {
	case indexTypeServerless:
		pcio.Fprintf(writer, "Cloud\t%s\n", options.cloud)
		pcio.Fprintf(writer, "Region\t%s\n", options.region)
		if options.sourceCollection != "" {
			pcio.Fprintf(writer, "Source Collection\t%s\n", options.sourceCollection)
		}
	case indexTypePod:
		pcio.Fprintf(writer, "Environment\t%s\n", options.environment)
		pcio.Fprintf(writer, "Pod Type\t%s\n", options.podType)
		pcio.Fprintf(writer, "Replicas\t%d\n", options.replicas)
		pcio.Fprintf(writer, "Shards\t%d\n", options.shards)
		if len(options.metadataConfig) > 0 {
			pcio.Fprintf(writer, "Metadata Config\t%s\n", text.InlineJSON(options.metadataConfig))
		}
		if options.sourceCollection != "" {
			pcio.Fprintf(writer, "Source Collection\t%s\n", options.sourceCollection)
		}
	case indexTypeIntegrated:
		pcio.Fprintf(writer, "Cloud\t%s\n", options.cloud)
		pcio.Fprintf(writer, "Region\t%s\n", options.region)
		pcio.Fprintf(writer, "Model\t%s\n", options.model)
		if len(options.fieldMap) > 0 {
			pcio.Fprintf(writer, "Field Map\t%s\n", text.InlineJSON(options.fieldMap))
		}
		if len(options.readParameters) > 0 {
			pcio.Fprintf(writer, "Read Parameters\t%s\n", text.InlineJSON(options.readParameters))
		}
		if len(options.writeParameters) > 0 {
			pcio.Fprintf(writer, "Write Parameters\t%s\n", text.InlineJSON(options.writeParameters))
		}
	}

	if len(options.tags) > 0 {
		pcio.Fprintf(writer, "\t\n")
		pcio.Fprintf(writer, "Tags\t%s\n", text.InlineJSON(options.tags))
	}

	writer.Flush()
	pcio.Println()
}

func renderSuccessOutput(idx *pinecone.Index, options createIndexOptions) {
	if options.json {
		json := text.IndentJSON(idx)
		pcio.Println(json)
		return
	}

	describeCommand := pcio.Sprintf("pc index describe %s", idx.Name)
	msg.SuccessMsg("Index %s created successfully. Run %s to check status. \n\n", style.Emphasis(idx.Name), style.Code(describeCommand))
	presenters.PrintDescribeIndexTable(idx)
}

// validate specific input params
func (c *createIndexOptions) validate() error {
	// Determine index type for validation
	idxType, err := c.deriveIndexType()
	if err != nil {
		return err
	}

	switch idxType {
	case indexTypeServerless:
		// Serverless requires cloud and region
		if c.cloud == "" {
			return pcio.Error("--cloud is required for serverless indexes")
		}
		if c.region == "" {
			return pcio.Error("--region is required for serverless indexes")
		}
		// Serverless cannot have pod-specific flags
		if c.environment != "" {
			return pcio.Error("--environment cannot be used with serverless indexes")
		}
		if c.podType != "" {
			return pcio.Error("--pod_type cannot be used with serverless indexes")
		}
		if c.shards != 1 {
			return pcio.Error("--shards cannot be used with serverless indexes")
		}
		if c.replicas != 1 {
			return pcio.Error("--replicas cannot be used with serverless indexes")
		}
		if len(c.metadataConfig) > 0 {
			return pcio.Error("--metadata_config cannot be used with serverless indexes")
		}

	case indexTypePod:
		// Pod requires environment
		if c.environment == "" {
			return pcio.Error("--environment is required for pod indexes")
		}
		// Pod requires pod_type
		if c.podType == "" {
			return pcio.Error("--pod_type is required for pod indexes")
		}
		// Pod cannot have serverless/integrated-specific flags
		if c.cloud != "" {
			return pcio.Error("--cloud cannot be used with pod indexes")
		}
		if c.region != "" {
			return pcio.Error("--region cannot be used with pod indexes")
		}
		if c.vectorType != "" {
			return pcio.Error("--vector_type cannot be used with pod indexes")
		}
		if c.model != "" {
			return pcio.Error("--model cannot be used with pod indexes")
		}
		if len(c.fieldMap) > 0 {
			return pcio.Error("--field_map cannot be used with pod indexes")
		}
		if len(c.readParameters) > 0 {
			return pcio.Error("--read_parameters cannot be used with pod indexes")
		}
		if len(c.writeParameters) > 0 {
			return pcio.Error("--write_parameters cannot be used with pod indexes")
		}

	case indexTypeIntegrated:
		// Integrated requires cloud and region
		if c.cloud == "" {
			return pcio.Error("--cloud is required for integrated indexes")
		}
		if c.region == "" {
			return pcio.Error("--region is required for integrated indexes")
		}
		// Integrated requires model
		if c.model == "" {
			return pcio.Error("--model is required for integrated indexes")
		}
		// Integrated cannot have pod-specific flags
		if c.environment != "" {
			return pcio.Error("--environment cannot be used with integrated indexes")
		}
		if c.podType != "" {
			return pcio.Error("--pod_type cannot be used with integrated indexes")
		}
		if c.shards != 1 {
			return pcio.Error("--shards cannot be used with integrated indexes")
		}
		if c.replicas != 1 {
			return pcio.Error("--replicas cannot be used with integrated indexes")
		}
		if len(c.metadataConfig) > 0 {
			return pcio.Error("--metadata_config cannot be used with integrated indexes")
		}
		if c.vectorType != "" {
			return pcio.Error("--vector_type cannot be used with integrated indexes")
		}
	}

	return nil
}

// determine the type of index being created based on high level input params
func (c *createIndexOptions) deriveIndexType() (indexType, error) {
	// Count how many index types are specified
	typeCount := 0
	if c.serverless {
		typeCount++
	}
	if c.pod {
		typeCount++
	}
	if c.integrated {
		typeCount++
	}

	// If multiple types are specified, that's an error
	if typeCount > 1 {
		return "", pcio.Error("only one index type can be specified. Use --serverless, --pod, or --integrated")
	}

	// If no type is explicitly specified, default to serverless
	if typeCount == 0 {
		c.serverless = true
	}

	// Determine type based on explicit flags
	if c.serverless {
		// Default to serverless index with common defaults
		if c.cloud == "" {
			c.cloud = "aws"
		}
		if c.region == "" {
			c.region = "us-east-1"
		}
		return indexTypeServerless, nil
	}

	if c.pod {
		return indexTypePod, nil
	}

	if c.integrated {
		return indexTypeIntegrated, nil
	}

	return "", pcio.Error("invalid index type")
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
