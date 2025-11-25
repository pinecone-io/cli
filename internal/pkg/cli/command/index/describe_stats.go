package index

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/flags"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

type describeStatsCmdOptions struct {
	indexName string
	filter    flags.JSONObject
	json      bool
}

func NewDescribeIndexStatsCmd() *cobra.Command {
	options := describeStatsCmdOptions{}
	cmd := &cobra.Command{
		Use:   "describe-stats",
		Short: "Describe the stats of an index",
		Long: help.Long(`
			Return index statistics including dimension, total vector count, namespaces summary, and metadata field counts.
			Use an optional metadata filter to restrict the scope of counts.

			JSON input may be inline, loaded from a file with @path, or read from stdin with @-.
		`),
		Example: help.Examples(`
			pc index describe-stats --index-name "index-name"
			pc index describe-stats --index-name "index-name" --filter '{"k":"v"}'
			pc index describe-stats --index-name "index-name" --filter @./filter.json
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runDescribeIndexStatsCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of index to describe stats for")
	cmd.Flags().VarP(&options.filter, "filter", "f", "metadata filter to apply to the operation (inline JSON, @path.json, or @- for stdin; max size: see PC_CLI_MAX_JSON_BYTES)")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	_ = cmd.MarkFlagRequired("index-name")

	return cmd
}

func runDescribeIndexStatsCmd(ctx context.Context, options describeStatsCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)

	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, "")
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	// Build metadata filter if provided
	var filter *pinecone.MetadataFilter
	if options.filter != nil {
		filter, err = pinecone.NewMetadataFilter(options.filter)
		if err != nil {
			msg.FailMsg("Failed to create filter: %s", err)
			exit.Errorf(err, "Failed to create filter")
		}
	}

	resp, err := ic.DescribeIndexStatsFiltered(ctx, filter)
	if err != nil {
		msg.FailMsg("Failed to describe stats: %s", err)
		exit.Error(err, "Failed to describe stats")
	}

	if options.json {
		json := text.IndentJSON(resp)
		pcio.Println(json)
	} else {
		presenters.PrintDescribeIndexStatsTable(resp)
	}
}
