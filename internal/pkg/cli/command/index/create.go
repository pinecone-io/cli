package index

import (
	"context"
	"fmt"
	"strconv"
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
	"github.com/spf13/pflag"
)

type indexType string

const (
	indexTypeServerless indexType = "serverless"
	indexTypeIntegrated indexType = "integrated"
	indexTypePod        indexType = "pod"
)

type createIndexOptions struct {
	// index name
	name string

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

	// confirmation
	yes bool

	// interactive mode
	interactive bool

	json bool
}

func NewCreateIndexCmd() *cobra.Command {
	options := createIndexOptions{}

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new index with the specified configuration",
		// No Args validation - we handle missing name in Run function
		Long: heredoc.Docf(`
		The %s command creates a new index with the specified configuration. By default, it creates a serverless index
		with sensible defaults. You can explicitly specify the index type using the appropriate flag:

			- Serverless (default): Use --serverless or no flag
			  %s
			- Pod: Use --pod flag
			  %s
			- Integrated: Use --integrated flag
			  %s

		Interactive Mode:
		Use --interactive or run without an index name to enter interactive mode, where you'll be prompted
		for each configuration value with sensible defaults.

		`, style.Code("pc index create"),
			style.URL(docslinks.DocsIndexCreate),
			style.URL(docslinks.DocsPodTypes),
			style.URL(docslinks.DocsIntegratedEmbedding)),
		Example: heredoc.Doc(`
		# create a serverless index (default)
		$ pc index create my-index

		# create a serverless index with custom configuration
		$ pc index create my-index --dimension 768 --metric euclidean --cloud aws --region us-east-1

		# create a pod index (with defaults)
		$ pc index create my-index --pod

		# create a pod index with custom configuration
		$ pc index create my-index --pod --environment us-east-1-aws --pod_type p1.x1 --shards 2 --replicas 2

		# create an integrated index (with defaults)
		$ pc index create my-index --integrated

		# create an integrated index with custom configuration
		$ pc index create my-index --integrated --cloud aws --region us-east-1 --model multilingual-e5-large --field_map text=chunk_text
		`),
		Run: func(cmd *cobra.Command, args []string) {
			var name string
			if len(args) > 0 {
				name = args[0]
			} else {
				// Automatically enable interactive mode if no name is provided
				options.interactive = true
			}
			runCreateIndexCmd(name, options, cmd)
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
	cmd.Flags().BoolVarP(&options.yes, "yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolVarP(&options.interactive, "interactive", "i", false, "Interactive mode - prompt for values")

	return cmd
}

func runCreateIndexCmd(name string, options createIndexOptions, cmd *cobra.Command) {
	ctx := context.Background()
	pc := sdk.NewPineconeClient()

	// Handle interactive mode
	if options.interactive || name == "" {
		options = runInteractiveMode(options)
		name = options.name // Get the name from interactive mode
	}

	// validate and derive index type from arguments
	err := options.validate(cmd)
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

	// Ask for confirmation unless --yes flag is used
	if !options.yes {
		if !confirmCreation(name) {
			msg.InfoMsg("Index creation cancelled.")
			return
		}
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

// Configuration map for interactive mode
// getRequiredFieldsForIndexType returns which fields should be prompted for in interactive mode
// based on the index type and vector type
func getRequiredFieldsForIndexType(idxType indexType, vectorType string) []string {
	switch idxType {
	case indexTypeServerless:
		fields := []string{"metric", "vectorType", "cloud", "region"}
		if vectorType == "dense" {
			fields = append(fields, "dimension")
		}
		return fields
	case indexTypePod:
		return []string{"dimension", "metric", "environment", "podType", "shards", "replicas"}
	case indexTypeIntegrated:
		fields := []string{"metric", "cloud", "region", "model", "fieldMap"}
		if vectorType == "dense" {
			fields = append(fields, "dimension")
		}
		return fields
	default:
		return []string{}
	}
}

// handleRegionPrompt prompts the user for region selection based on index type and cloud
func handleRegionPrompt(options *createIndexOptions, idxType indexType) {
	// Determine which region map to use based on index type
	var regionsMap map[string][]string
	if idxType == indexTypePod {
		regionsMap = podCloudRegions
	} else {
		// Serverless or integrated
		regionsMap = serverlessCloudRegions
	}

	// Show available regions for the selected cloud
	if regions, exists := regionsMap[options.cloud]; exists && len(regions) > 0 {
		defaultRegion := regions[0] // First region is the default

		pcio.Printf("Available regions for %s:\n", options.cloud)
		for i, region := range regions {
			if region == defaultRegion {
				pcio.Printf("  %d. %s (default)\n", i+1, region)
			} else {
				pcio.Printf("  %d. %s\n", i+1, region)
			}
		}
		pcio.Printf("Enter region [%s]: ", defaultRegion)

		var response string
		fmt.Scanln(&response)

		if response == "" {
			options.region = defaultRegion
		} else {
			// Check if user entered a number (choice) or region name
			if choiceIndex, err := strconv.Atoi(response); err == nil && choiceIndex > 0 && choiceIndex <= len(regions) {
				options.region = regions[choiceIndex-1]
			} else {
				options.region = response
			}
		}
	} else {
		// Fallback if no regions found for the cloud
		pcio.Printf("Region: ")
		fmt.Scanln(&options.region)
	}
}

// Valid regions for serverless/integrated indexes for each cloud provider
var serverlessCloudRegions = map[string][]string{
	"aws":   {"us-east-1", "us-west-2", "eu-west-1"},
	"gcp":   {"us-central1", "europe-west4"},
	"azure": {"eastus2"},
}

// Valid regions for pod indexes for each cloud provider
var podCloudRegions = map[string][]string{
	"aws":   {"us-east-1"},
	"gcp":   {"us-west1-gcp", "us-central1-gcp", "us-west4-gcp", "us-east4-gcp", "northamerica-northeast1-gcp", "asia-northeast1-gcp", "asia-southeast1-gcp", "us-east1-gcp", "eu-west1-gcp", "eu-west4-gcp"},
	"azure": {"eastus-azure"},
}

// Available embedding models for integrated indexes
var denseEmbeddingModels = []string{
	"multilingual-e5-large",
	"llama-text-embed-v2",
}

var sparseEmbeddingModels = []string{
	"pinecone-sparse-english-v0",
}

// runInteractiveMode prompts the user for index configuration values
func runInteractiveMode(options createIndexOptions) createIndexOptions {
	pcio.Println("===================================")
	pcio.Println("| Interactive index creation mode |")
	pcio.Println("===================================")
	pcio.Println()

	// Get index name
	options.name = promptForString("Index name", "")

	// Get index type
	indexType := promptForChoice("Index type", []string{"serverless", "pod", "integrated"}, "serverless")
	switch indexType {
	case "serverless":
		options.serverless = true
	case "pod":
		options.pod = true
	case "integrated":
		options.integrated = true
	}

	// Get vector type for serverless and integrated
	var vectorType string
	if options.serverless || options.integrated {
		vectorType = promptForChoice("Vector type", []string{"dense", "sparse"}, "dense")
		options.vectorType = vectorType
	}

	// Determine index type for field requirements
	idxType, _ := options.deriveIndexType()
	requiredFields := getRequiredFieldsForIndexType(idxType, vectorType)

	// Prompt for each required field
	for _, field := range requiredFields {
		switch field {
		case "dimension":
			options.dimension = int32(promptForInt("Dimension", int(options.dimension)))
		case "metric":
			options.metric = promptForChoice("Metric", []string{"cosine", "euclidean", "dotproduct"}, options.metric)
		case "vectorType":
			// Already handled above
		case "cloud":
			options.cloud = promptForChoice("Cloud provider", []string{"aws", "gcp", "azure"}, "aws")
		case "region":
			// Handle region prompting with cloud-specific options
			handleRegionPrompt(&options, idxType)
		case "environment":
			// For pod indexes, we need to prompt for environment (which includes cloud and region)
			// Show available environments based on the documentation
			environments := []string{
				"us-west1-gcp", "us-central1-gcp", "us-west4-gcp", "us-east4-gcp",
				"northamerica-northeast1-gcp", "asia-northeast1-gcp", "asia-southeast1-gcp",
				"us-east1-gcp", "eu-west1-gcp", "eu-west4-gcp", "us-east-1-aws", "eastus-azure",
			}
			options.environment = promptForChoice("Environment", environments, "us-east-1-aws")
		case "podType":
			options.podType = promptForChoice("Pod type", []string{"p1.x1", "p1.x2", "p1.x4", "p1.x8", "s1.x1", "s1.x2", "s1.x4", "s1.x8", "p2.x1", "p2.x2", "p2.x4", "p2.x8"}, "p1.x1")
		case "shards":
			options.shards = int32(promptForInt("Shards", int(options.shards)))
		case "replicas":
			options.replicas = int32(promptForInt("Replicas", int(options.replicas)))
		case "model":
			// Show appropriate models based on vector type
			var models []string
			var defaultModel string
			if vectorType == "dense" {
				models = denseEmbeddingModels
				defaultModel = "multilingual-e5-large"
			} else {
				models = sparseEmbeddingModels
				defaultModel = "pinecone-sparse-english-v0"
			}
			options.model = promptForChoice("Model", models, defaultModel)
		case "fieldMap":
			// For now, use a simple default field map
			options.fieldMap = map[string]string{"text": "chunk_text"}
		}
	}

	pcio.Println()
	return options
}

// promptForString prompts the user for a string value
func promptForString(prompt string, defaultValue string) string {
	if defaultValue != "" {
		pcio.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		pcio.Printf("%s: ", prompt)
	}

	var response string
	fmt.Scanln(&response)

	if response == "" && defaultValue != "" {
		return defaultValue
	}
	return response
}

// promptForInt prompts the user for an integer value
func promptForInt(prompt string, defaultValue int) int {
	pcio.Printf("%s [%d]: ", prompt, defaultValue)

	var response string
	fmt.Scanln(&response)

	if response == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(response)
	if err != nil {
		pcio.Printf("Invalid number, using default %d\n", defaultValue)
		return defaultValue
	}
	return val
}

// promptForChoice prompts the user to choose from a list of options
func promptForChoice(prompt string, choices []string, defaultValue string) string {
	pcio.Printf("%s:\n", prompt)
	for i, choice := range choices {
		if choice == defaultValue {
			pcio.Printf("  %d. %s (default)\n", i+1, choice)
		} else {
			pcio.Printf("  %d. %s\n", i+1, choice)
		}
	}

	pcio.Printf("Enter choice [%s]: ", defaultValue)
	var response string
	fmt.Scanln(&response)

	if response == "" {
		return defaultValue
	}

	choiceIndex, err := strconv.Atoi(response)
	if err != nil || choiceIndex < 1 || choiceIndex > len(choices) {
		pcio.Printf("Invalid choice, using default %s\n", defaultValue)
		return defaultValue
	}

	return choices[choiceIndex-1]
}

// confirmCreation prompts the user for confirmation to create the index
func confirmCreation(name string) bool {
	pcio.Printf("Create index '%s'? [y/N]: ", style.Emphasis(name))

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		// If there's an error reading input, assume no
		return false
	}

	// Convert to lowercase and trim whitespace
	response = strings.ToLower(strings.TrimSpace(response))

	// Accept y, yes, Y, YES
	return response == "y" || response == "yes"
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
func (c *createIndexOptions) validate(cmd *cobra.Command) error {
	// Determine index type for validation
	idxType, err := c.deriveIndexType()
	if err != nil {
		return err
	}

	// Define which flags are invalid for each index type
	invalidFlags := map[indexType]map[string]string{
		indexTypeServerless: {
			"environment":      "--environment cannot be used with serverless indexes",
			"pod_type":         "--pod_type cannot be used with serverless indexes",
			"metadata_config":  "--metadata_config cannot be used with serverless indexes",
			"shards":           "--shards cannot be used with serverless indexes",
			"replicas":         "--replicas cannot be used with serverless indexes",
			"model":            "--model cannot be used with serverless indexes",
			"field_map":        "--field_map cannot be used with serverless indexes",
			"read_parameters":  "--read_parameters cannot be used with serverless indexes",
			"write_parameters": "--write_parameters cannot be used with serverless indexes",
		},
		indexTypePod: {
			"cloud":            "--cloud cannot be used with pod indexes",
			"region":           "--region cannot be used with pod indexes",
			"vector_type":      "--vector_type cannot be used with pod indexes",
			"model":            "--model cannot be used with pod indexes",
			"field_map":        "--field_map cannot be used with pod indexes",
			"read_parameters":  "--read_parameters cannot be used with pod indexes",
			"write_parameters": "--write_parameters cannot be used with pod indexes",
		},
		indexTypeIntegrated: {
			"environment":     "--environment cannot be used with integrated indexes",
			"pod_type":        "--pod_type cannot be used with integrated indexes",
			"metadata_config": "--metadata_config cannot be used with integrated indexes",
			"shards":          "--shards cannot be used with integrated indexes",
			"replicas":        "--replicas cannot be used with integrated indexes",
		},
	}

	// Check only flags that were explicitly set by the user
	var validationError error
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if errorMsg, exists := invalidFlags[idxType][flag.Name]; exists {
			validationError = pcio.Error(errorMsg)
		}
	})

	// Additional validation for integrated indexes
	if idxType == indexTypeIntegrated && c.model != "" {
		// Validate model based on vector type
		var validModels []string
		if c.vectorType == "dense" {
			validModels = denseEmbeddingModels
		} else if c.vectorType == "sparse" {
			validModels = sparseEmbeddingModels
		}

		// Check if the model is valid
		modelValid := false
		for _, validModel := range validModels {
			if c.model == validModel {
				modelValid = true
				break
			}
		}

		if !modelValid {
			validationError = pcio.Error(fmt.Sprintf("invalid model '%s' for vector type '%s'. Valid models are: %s", c.model, c.vectorType, strings.Join(validModels, ", ")))
		}
	}

	// Handle dimension logic for different vector types
	if c.vectorType == "sparse" {
		// Check if dimension flag was explicitly set by user
		if cmd.Flags().Changed("dimension") {
			validationError = pcio.Error("--dimension cannot be used with sparse vector type. Sparse vectors have variable dimensions determined at runtime.")
		} else {
			// For sparse vectors without explicit dimension, set to 0 so pointerOrNil returns nil
			c.dimension = 0
		}

		// Check that sparse vectors use dotproduct metric
		if c.metric != "" && c.metric != "dotproduct" {
			validationError = pcio.Error("sparse vector type requires dotproduct metric")
		}
	}

	return validationError
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
		// Default to a common environment if none specified
		if c.environment == "" {
			c.environment = "us-east-1-aws"
		}
		// Default to a common pod type if none specified
		if c.podType == "" {
			c.podType = "p1.x1"
		}
		return indexTypePod, nil
	}

	if c.integrated {
		// Default to common cloud and region if none specified
		if c.cloud == "" {
			c.cloud = "aws"
		}
		if c.region == "" {
			c.region = "us-east-1"
		}
		// Default to a common model if none specified
		if c.model == "" {
			c.model = "multilingual-e5-large"
		}
		// Default to a common field_map if none specified
		if len(c.fieldMap) == 0 {
			c.fieldMap = map[string]string{"text": "chunk_text"}
		}
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
