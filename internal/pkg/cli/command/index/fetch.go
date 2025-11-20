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

type fetchCmdOptions struct {
	name            string
	namespace       string
	ids             []string
	filter          flags.JSONObject
	limit           uint32
	paginationToken string
	json            bool
}

func NewFetchCmd() *cobra.Command {
	options := fetchCmdOptions{}
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch vectors by ID or metadata filter from an index",
		Example: help.Examples(`
			pc index fetch --name my-index --ids 123,456,789
			pc index fetch --name my-index --filter '{"key": "value"}'
			pc index fetch --name my-index --filter @./filter.json
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runFetchCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringSliceVarP(&options.ids, "ids", "i", []string{}, "IDs of vectors to fetch")
	cmd.Flags().VarP(&options.filter, "filter", "f", "metadata filter to apply to the fetch")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the index to fetch from")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to fetch from")
	cmd.Flags().Uint32VarP(&options.limit, "limit", "l", 0, "maximum number of vectors to fetch")
	cmd.Flags().StringVarP(&options.paginationToken, "pagination-token", "p", "", "pagination token to continue a previous listing operation")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")

	cmd.MarkFlagsMutuallyExclusive("ids", "filter")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runFetchCmd(ctx context.Context, options fetchCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)

	// Default namespace
	ns := options.namespace
	if ns == "" {
		ns = "__default__"
	}

	if len(options.ids) > 0 && (options.limit > 0 || options.paginationToken != "") {
		msg.FailMsg("ids and limit/pagination-token cannot be used together")
		exit.ErrorMsg("ids and limit/pagination-token cannot be used together")
	}

	ic, err := sdk.NewIndexConnection(ctx, pc, options.name, ns)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	// Fetch vectors by ID
	if len(options.ids) > 0 {
		vectors, err := ic.FetchVectors(ctx, options.ids)
		if err != nil {
			exit.Error(err, "Failed to fetch vectors")
		}
		printFetchVectorsResults(presenters.NewFetchVectorsResultsFromFetch(vectors), options)
	}

	// Fetch vectors by metadata filter
	if options.filter != nil {
		filter, err := pinecone.NewMetadataFilter(options.filter)
		if err != nil {
			msg.FailMsg("Failed to create filter: %s", err)
			exit.Errorf(err, "Failed to create filter")
		}

		req := &pinecone.FetchVectorsByMetadataRequest{
			Filter: filter,
		}

		if options.limit > 0 {
			req.Limit = &options.limit
		}
		if options.paginationToken != "" {
			req.PaginationToken = &options.paginationToken
		}

		vectors, err := ic.FetchVectorsByMetadata(ctx, req)
		if err != nil {
			exit.Error(err, "Failed to fetch vectors by metadata")
		}
		printFetchVectorsResults(presenters.NewFetchVectorsResultsFromFetchByMetadata(vectors), options)
	}
}

func printFetchVectorsResults(results *presenters.FetchVectorsResults, options fetchCmdOptions) {
	if options.json {
		json := text.IndentJSON(results)
		pcio.Println(json)
	} else {
		presenters.PrintFetchVectorsTable(results)
	}
}
