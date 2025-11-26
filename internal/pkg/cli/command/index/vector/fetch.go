package vector

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/argio"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/flags"
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

// FetchBody is the JSON payload schema for --body.
// Fields: ids, filter, limit, pagination_token.
// When ids are provided, pagination fields are not applicable.
type FetchBody struct {
	Ids             []string       `json:"ids"`
	Filter          map[string]any `json:"filter"`
	Limit           *uint32        `json:"limit"`
	PaginationToken *string        `json:"pagination_token"`
}

type fetchCmdOptions struct {
	indexName       string
	namespace       string
	ids             flags.StringList
	filter          flags.JSONObject
	limit           uint32
	paginationToken string
	body            string
	json            bool
}

func NewFetchCmd() *cobra.Command {
	options := fetchCmdOptions{}
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch vectors by ID or metadata filter from an index",
		Long: help.Long(`
			Fetch vectors from an index either by explicit IDs or by a metadata filter with optional pagination.

			When using --ids, pagination flags (--limit, --pagination-token) are not applicable.
			JSON inputs may be inline, loaded from ./file.json[l], or read from stdin with '-'.
			A --body payload can supply ids, filter, limit, and pagination_token fields. Flags win if both are provided.
		`),
		Example: help.Examples(`
			pc index vector fetch --index-name my-index --ids '["123","456","789"]'
			pc index vector fetch --index-name my-index --ids ./ids.json

			pc index vector fetch --index-name my-index --filter '{"genre":{"$eq":"rock"}}'
			pc index vector fetch --index-name my-index --filter ./filter.json

			pc index vector fetch --index-name my-index --body ./fetch.json
			cat fetch.json | pc index vector fetch --index-name my-index --body -
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runFetchCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().VarP(&options.ids, "ids", "i", "IDs of vectors to fetch (inline JSON array, ./path.json, or '-' for stdin)")
	cmd.Flags().VarP(&options.filter, "filter", "f", "metadata filter to apply to the fetch (inline JSON, ./path.json, or '-' for stdin)")
	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to fetch from")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to fetch from")
	cmd.Flags().Uint32VarP(&options.limit, "limit", "l", 0, "maximum number of vectors to fetch")
	cmd.Flags().StringVarP(&options.paginationToken, "pagination-token", "p", "", "pagination token to continue a previous listing operation")
	cmd.Flags().StringVar(&options.body, "body", "", "request body JSON (inline, ./path.json, or '-' for stdin; only one argument may use stdin)")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")

	cmd.MarkFlagsMutuallyExclusive("ids", "filter")
	_ = cmd.MarkFlagRequired("index-name")

	return cmd
}

func runFetchCmd(ctx context.Context, options fetchCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)

	// Apply body overlay if provided
	if options.body != "" {
		if b, src, err := argio.DecodeJSONArg[FetchBody](options.body); err != nil {
			msg.FailMsg("Failed to parse fetch body (%s): %s", style.Emphasis(src.Label), err)
			exit.Errorf(err, "Failed to parse fetch body (%s): %v", src.Label, err)
		} else if b != nil {
			if len(options.ids) == 0 && len(b.Ids) > 0 {
				options.ids = b.Ids
			}
			if options.filter == nil && b.Filter != nil {
				options.filter = b.Filter
			}
			if b.Limit != nil {
				options.limit = *b.Limit
			}
			if b.PaginationToken != nil {
				options.paginationToken = *b.PaginationToken
			}
		}
	}

	if len(options.ids) > 0 && (options.limit > 0 || options.paginationToken != "") {
		msg.FailMsg("ids and limit/pagination-token cannot be used together")
		exit.ErrorMsg("ids and limit/pagination-token cannot be used together")
	}

	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, options.namespace)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	if options.ids == nil && options.filter == nil {
		msg.FailMsg("Either --ids or --filter must be provided")
		exit.ErrorMsg("Either --ids or --filter must be provided")
	}

	// Fetch vectors by ID
	if len(options.ids) > 0 {
		vectors, err := ic.FetchVectors(ctx, []string(options.ids))
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
