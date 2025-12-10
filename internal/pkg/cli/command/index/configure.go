package index

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

type configureIndexOptions struct {
	// required for index lookup
	name string

	// pods
	podType  string
	replicas int32

	// integrated
	model           string
	fieldMap        map[string]string
	readParameters  map[string]string
	writeParameters map[string]string

	// serverless & integrated
	readMode     string
	readNodeType string
	readShards   int32
	readReplicas int32

	// optional for all index types
	deletionProtection string
	tags               map[string]string

	json bool
}

func NewConfigureIndexCmd() *cobra.Command {
	options := configureIndexOptions{}

	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure an existing index",
		Example: help.Examples(`
			pc index configure --name "index-name" --deletion-protection "enabled"
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runConfigureIndexCmd(cmd.Context(), cmd, options)
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "Name of index to configure")

	// pods
	cmd.Flags().StringVarP(&options.podType, "pod-type", "t", "", "Type of pod to use, can only upgrade when configuring")
	cmd.Flags().Int32VarP(&options.replicas, "replicas", "r", 0, "Replicas of the index to configure")

	// integrated
	cmd.Flags().StringVar(&options.model, "model", "", "The name of the embedding model to use for the index")
	cmd.Flags().StringToStringVar(&options.fieldMap, "field-map", map[string]string{}, "Identifies the name of the text field from your document model that will be embedded")
	cmd.Flags().StringToStringVar(&options.readParameters, "read-parameters", map[string]string{}, "The read parameters for the embedding model")
	cmd.Flags().StringToStringVar(&options.writeParameters, "write-parameters", map[string]string{}, "The write parameters for the embedding model")

	// serverless & integrated
	cmd.Flags().StringVar(&options.readMode, "read-mode", "", "The read capacity mode to use. One of: ondemand, dedicated. When converting from ondemand to dedicated, you must provide read-node-type, read-shards, and read-replicas")
	cmd.Flags().StringVar(&options.readNodeType, "read-node-type", "", "The type of machines to use. Available options: b1 and t1. t1 includes increased processing power and memory")
	cmd.Flags().Int32Var(&options.readShards, "read-shards", 0, "The number of shards to use. Shards determine the storage capacity of an index, with each shard providing 250 GB of storage")
	cmd.Flags().Int32Var(&options.readReplicas, "read-replicas", 0, "The number of replicas to use. Replicas duplicate the compute resources and data of an index, allowing higher query throughput and availability")

	// optional for all index types
	cmd.Flags().StringVarP(&options.deletionProtection, "deletion-protection", "p", "", "Enable or disable deletion protection for the index. One of: enabled, disabled")
	cmd.Flags().StringToStringVar(&options.tags, "tags", map[string]string{}, "Custom user tags to add to an index")

	return cmd
}

func runConfigureIndexCmd(ctx context.Context, cmd *cobra.Command, options configureIndexOptions) {
	pc := sdk.NewPineconeClient(ctx)

	// index tags
	var indexTags pinecone.IndexTags
	if len(options.tags) > 0 {
		indexTags = pinecone.IndexTags(options.tags)
	}

	// embed configuration
	var embed *pinecone.ConfigureIndexEmbed
	if options.model != "" || len(options.fieldMap) > 0 || len(options.readParameters) > 0 || len(options.writeParameters) > 0 {
		fieldMap := toInterfaceMap(options.fieldMap)
		readParameters := toInterfaceMap(options.readParameters)
		writeParameters := toInterfaceMap(options.writeParameters)

		embed = &pinecone.ConfigureIndexEmbed{}

		if options.model != "" {
			embed.Model = &options.model
		}
		if fieldMap != nil {
			embed.FieldMap = &fieldMap
		}
		if readParameters != nil {
			embed.ReadParameters = &readParameters
		}
		if writeParameters != nil {
			embed.WriteParameters = &writeParameters
		}
	}

	// read capacity configuration
	var mode *string
	var nodeType *string
	var shards *int32
	var replicas *int32
	if cmd.Flags().Changed("read-mode") {
		mode = &options.readMode
	}
	if cmd.Flags().Changed("read-node-type") {
		nodeType = &options.readNodeType
	}
	if cmd.Flags().Changed("read-shards") {
		shards = &options.readShards
	}
	if cmd.Flags().Changed("read-replicas") {
		replicas = &options.readReplicas
	}
	readCapacity, err := constructReadCapacity(mode, nodeType, shards, replicas)
	if err != nil {
		msg.FailMsg("Failed to configure index %s: %+v\n", style.Emphasis(options.name), err)
		exit.Error(err, "Failed to configure index")
	}

	idx, err := pc.ConfigureIndex(ctx, options.name, pinecone.ConfigureIndexParams{
		PodType:            options.podType,
		Replicas:           options.replicas,
		DeletionProtection: pinecone.DeletionProtection(options.deletionProtection),
		ReadCapacity:       readCapacity,
		Tags:               indexTags,
		Embed:              embed,
	})
	if err != nil {
		msg.FailMsg("Failed to configure index %s: %+v\n", style.Emphasis(options.name), err)
		exit.Error(err, "Failed to configure index")
	}
	if options.json {
		json := text.IndentJSON(idx)
		pcio.Println(json)
		return
	}

	describeCommand := pcio.Sprintf("pc index describe --name %s", idx.Name)
	msg.SuccessMsg("Index %s configured successfully. Run %s to check status. \n\n", style.Emphasis(idx.Name), style.Code(describeCommand))
	presenters.PrintDescribeIndexTable(idx)
}
