package index

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/docslinks"
	errorutil "github.com/pinecone-io/cli/internal/pkg/utils/error"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/index"
	indexpresenters "github.com/pinecone-io/cli/internal/pkg/utils/index/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/interactive"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/models"
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
	interactive   bool
}

func NewCreateIndexCmd() *cobra.Command {
	options := createIndexOptions{}

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new index with the specified configuration",
		Long: heredoc.Docf(`
		The %s command creates a new index with the specified configuration. There are different types of indexes
		you can create:

			- Serverless (dense or sparse)
			- Pod (dense only)

		For serverless indexes, you can specify an embedding model to use via the %s flag:

		The CLI will try to automatically infer missing settings from those provided.

		Use the %s flag to enable interactive mode, which will guide you through configuring
		the index settings step by step.

		For detailed documentation, see:
		%s
		`, style.Code("pc index create"),
			style.Emphasis("--model"),
			style.Code("--interactive"),
			style.URL(docslinks.DocsIndexCreate)),
		Example: heredoc.Doc(`
		# create default index (serverless)
		$ pc index create my-index

		# create serverless index
		$ pc index create my-index --serverless

		# create pod index
		$ pc index create my-index --pod	

		# create a serverless index with explicit model
		$ pc index create my-index --model llama-text-embed-v2 --cloud aws --region us-east-1

		# create a serverless index with the default dense model
		$ pc index create my-index --model dense --cloud aws --region us-east-1

		# create a serverless index with the default sparse model
		$ pc index create my-index --model sparse --cloud aws --region us-east-1

		# create an index using interactive mode
		$ pc index create my-index --interactive

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
	cmd.Flags().StringVar(&options.CreateOptions.Model.Value, "model", "", fmt.Sprintf("Embedding model to use (e.g., llama-text-embed-v2, default, sparse). Use %s to see available models", style.Code("pc models")))
	cmd.Flags().StringToStringVar(&options.CreateOptions.FieldMap.Value, "field_map", map[string]string{}, "Identifies the name of the text field from your document model that will be embedded")
	cmd.Flags().StringToStringVar(&options.CreateOptions.ReadParameters.Value, "read_parameters", map[string]string{}, "The read parameters for the embedding model")
	cmd.Flags().StringToStringVar(&options.CreateOptions.WriteParameters.Value, "write_parameters", map[string]string{}, "The write parameters for the embedding model")

	// Optional flags
	cmd.Flags().Int32VarP(&options.CreateOptions.Dimension.Value, "dimension", "d", 0, "Dimension of the index to create")
	cmd.Flags().StringVarP(&options.CreateOptions.Metric.Value, "metric", "m", "", "Metric to use. One of: cosine, euclidean, dotproduct")
	cmd.Flags().StringVar(&options.CreateOptions.DeletionProtection.Value, "deletion_protection", "", "Whether to enable deletion protection for the index. One of: enabled, disabled")
	cmd.Flags().StringToStringVar(&options.CreateOptions.Tags.Value, "tags", map[string]string{}, "Custom user tags to add to an index")

	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")
	cmd.Flags().BoolVarP(&options.interactive, "interactive", "i", false, "Enable interactive mode to configure index settings step by step")

	return cmd
}

func collectInteractiveConfiguration(ctx context.Context, options index.CreateOptions) (index.CreateOptions, bool) {
	pcio.Println(style.Hint("Press Esc or Ctrl+C at any time to exit interactive mode.\n"))

	// Variables for model data (will be populated when needed)
	var availableModels []models.ModelInfo
	var modelsErr error

	// Index type selection - determine default based on existing flags
	var defaultChoice string
	if options.Serverless.Value {
		defaultChoice = "Serverless"
	} else if options.Pod.Value {
		defaultChoice = "Pod"
	} else {
		defaultChoice = "Serverless"
	}

	choice, exit := interactive.GetChoice(
		"Select index type",
		[]string{
			"Serverless",
			"Pod",
		},
		defaultChoice)

	if exit {
		return options, true
	}

	switch choice {
	case "Serverless":
		options.Serverless.Value = true
		options.Pod.Value = false
	case "Pod":
		options.Serverless.Value = false
		options.Pod.Value = true
	}

	// Serverless configuration
	if options.Serverless.Value {
		// Fetch available models for serverless
		availableModels, modelsErr = models.GetModels(ctx, true)
		if modelsErr != nil {
			pcio.Println(style.WarnMsg("Warning: Could not fetch available models!"))
		}

		// Model selection
		if modelsErr != nil {
			options.Model.Value = "llama-text-embed-v2"
		} else {
			// Create model choices
			modelChoices := make([]string, 0, len(availableModels)+1)
			modelChoices = append(modelChoices, "None (custom vectors)")

			for _, model := range availableModels {
				modelChoices = append(modelChoices, model.Model)
			}

			// Determine default model choice
			var defaultModelChoice string
			if options.Model.Value != "" {
				defaultModelChoice = options.Model.Value
			} else {
				defaultModelChoice = "llama-text-embed-v2" // Default dense model
			}

			modelChoice, exit := interactive.GetChoice("Select inference model", modelChoices, defaultModelChoice)
			if exit {
				return options, true
			}

			if modelChoice == "None (custom vectors)" {
				options.Model.Value = ""
			} else {
				options.Model.Value = modelChoice
			}
		}

		// Cloud and region
		cloud, exit := interactive.GetInput("Cloud provider (aws, gcp, azure)", options.Cloud.Value)
		if exit {
			return options, true
		}
		options.Cloud.Value = cloud

		region, exit := interactive.GetInput("Region (e.g., us-east-1)", options.Region.Value)
		if exit {
			return options, true
		}
		options.Region.Value = region

		// Vector type (only for serverless without model)
		if options.Model.Value == "" {
			vectorType, exit := interactive.GetChoice("Vector type", []string{"dense", "sparse"}, options.VectorType.Value)
			if exit {
				return options, true
			}
			options.VectorType.Value = vectorType
		}
	}

	// Environment (for pod)
	if options.Pod.Value {
		environment, exit := interactive.GetInput("Environment", options.Environment.Value)
		if exit {
			return options, true
		}
		options.Environment.Value = environment

		podType, exit := interactive.GetInput("Pod type", options.PodType.Value)
		if exit {
			return options, true
		}
		options.PodType.Value = podType

		shards, exit := interactive.GetIntInput("Number of shards", int(options.Shards.Value))
		if exit {
			return options, true
		}
		options.Shards.Value = int32(shards)

		replicas, exit := interactive.GetIntInput("Number of replicas", int(options.Replicas.Value))
		if exit {
			return options, true
		}
		options.Replicas.Value = int32(replicas)
	}

	// Common settings

	// Handle dimension based on vector type
	// Sparse models always use dimension 0, dense models may support multiple dimensions
	isSparse := false
	if options.Model.Value != "" {
		// Check if the selected model is sparse
		if modelsErr == nil {
			for _, model := range availableModels {
				if model.Model == options.Model.Value && model.VectorType != nil && *model.VectorType == "sparse" {
					isSparse = true
					break
				}
			}
		}
	} else {
		// For custom vectors, check the vector type
		isSparse = options.VectorType.Value == "sparse"
	}

	if isSparse {
		// Sparse vectors always use dimension 0
		options.Dimension.Value = 0
	} else {
		// Ask for dimension for dense models and custom vectors
		dimension, exit := interactive.GetIntInput("Dimension (0 for auto)", int(options.Dimension.Value))
		if exit {
			return options, true
		}
		options.Dimension.Value = int32(dimension)
	}

	// Only ask for metric if not sparse (sparse always uses dotproduct)
	// isSparse is already determined above

	if !isSparse {
		metric, exit := interactive.GetChoice("Metric", []string{"cosine", "euclidean", "dotproduct"}, options.Metric.Value)
		if exit {
			return options, true
		}
		options.Metric.Value = metric
	} else {
		// Sparse vectors always use dotproduct
		options.Metric.Value = "dotproduct"
	}

	// Set default deletion protection to disabled if not already set
	defaultDeletionProtection := options.DeletionProtection.Value
	if defaultDeletionProtection == "" {
		defaultDeletionProtection = "disabled"
	}

	deletionProtection, exit := interactive.GetChoice("Deletion protection", []string{"enabled", "disabled"}, defaultDeletionProtection)
	if exit {
		return options, true
	}
	options.DeletionProtection.Value = deletionProtection

	// Tags
	useTags := interactive.GetConfirmation("Add custom tags?")

	if useTags {
		// Initialize tags map if it doesn't exist
		if options.Tags.Value == nil {
			options.Tags.Value = make(map[string]string)
		}

		pcio.Println("Enter tags in key=value format. Press Enter with empty input to finish.")

		for {
			// Show current tags
			if len(options.Tags.Value) > 0 {
				for k, v := range options.Tags.Value {
					pcio.Printf("  %s=%s\n", style.Emphasis(k), style.ResourceName(v))
				}
				pcio.Println()
			}

			tagInput, exit := interactive.GetInput("Tag (key=value)", "")
			if exit {
				return options, true
			}

			// Empty input means done adding tags
			if strings.TrimSpace(tagInput) == "" {
				break
			}

			// Parse key=value format
			parts := strings.SplitN(tagInput, "=", 2)
			if len(parts) != 2 {
				pcio.Println(style.FailMsg("Invalid format. Please use key=value format."))
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if key == "" || value == "" {
				pcio.Println(style.FailMsg("Both key and value must be non-empty."))
				continue
			}

			// Add the tag
			options.Tags.Value[key] = value
		}
	}

	return options, false
}

func runCreateIndexCmd(options createIndexOptions, cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	// If interactive mode is enabled, collect configuration interactively
	if options.interactive {
		var exit bool
		options.CreateOptions, exit = collectInteractiveConfiguration(ctx, options.CreateOptions)
		if exit {
			pcio.Println(style.InfoMsg("Interactive mode cancelled."))
			return
		}
	}

	// validationErrors := index.ValidateCreateOptions(options.CreateOptions)
	// if len(validationErrors) > 0 {
	// 	msg.FailMsgMultiLine(validationErrors...)
	// 	exit.Error(errors.New(validationErrors[0])) // Use first error for exit code
	// }

	inferredOptions := index.InferredCreateOptions(ctx, options.CreateOptions)
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

	// Ask for user confirmation unless -y flag is set
	assumeYes, _ := cmd.Flags().GetBool("assume-yes")
	if !assumeYes {
		question := "Is this configuration correct? Do you want to proceed with creating the index?"
		if !interactive.GetConfirmation(question) {
			pcio.Println(style.InfoMsg("Index creation cancelled."))
			return
		}
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
			errorutil.HandleAPIError(err, cmd, args)
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
			errorutil.HandleAPIError(err, cmd, args)
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
			errorutil.HandleAPIError(err, cmd, args)
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
