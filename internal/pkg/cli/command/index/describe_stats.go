package index

import (
	"context"
	"encoding/json"
	"os"

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

type describeStatsCmdOptions struct {
	name       string
	filter     string
	filterFile string
	json       bool
}

func NewDescribeIndexStatsCmd() *cobra.Command {
	options := describeStatsCmdOptions{}
	cmd := &cobra.Command{
		Use:   "describe-stats",
		Short: "Describe the stats of an index",
		Example: help.Examples(`
			pc index describe-stats --name "index-name"
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runDescribeIndexStatsCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to describe stats for")
	cmd.Flags().StringVar(&options.filter, "filter", "", "filter to apply to the stats")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runDescribeIndexStatsCmd(ctx context.Context, options describeStatsCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)

	ic, err := sdk.NewIndexConnection(ctx, pc, options.name, "")
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	// Build metadata filter if provided
	var filter *pinecone.MetadataFilter
	if options.filter != "" || options.filterFile != "" {
		if options.filterFile != "" {
			raw, err := os.ReadFile(options.filterFile)
			if err != nil {
				msg.FailMsg("Failed to read filter file %s: %s", style.Emphasis(options.filterFile), err)
				exit.Errorf(err, "Failed to read filter file %s", options.filterFile)
			}
			options.filter = string(raw)
		}

		var filterMap map[string]any
		if err := json.Unmarshal([]byte(options.filter), &filterMap); err != nil {
			msg.FailMsg("Failed to parse filter: %s", err)
			exit.Errorf(err, "Failed to parse filter")
		}
		filter, err = pinecone.NewMetadataFilter(filterMap)
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
