package record

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
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

const defaultSearchTopK = 10

type searchCmdOptions struct {
	indexName string
	namespace string
	topK      int
	inputs    flags.JSONObject
	id        string
	fields    flags.StringList
	body      string
	json      bool
}

func NewSearchCmd() *cobra.Command {
	options := searchCmdOptions{
		topK: defaultSearchTopK,
	}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search records in an integrated index",
		Long: help.Long(`
			Run semantic search against records in an integrated index.

			Provide query text via --inputs (inline JSON, ./path.json, or '-' for stdin).
			You may also supply a full request body with --body to access advanced
			parameters like reranking, vector overrides, filters, or match_terms.
		`),
		Example: help.Examples(`
			pc index record search --index-name my-index --namespace my-namespace --inputs '{"text":"find similar"}'
			echo '{"text":"disease prevention"}' | pc index record search --index-name my-index --inputs -
			pc index record search --index-name my-index --namespace my-namespace --body ./search.json
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runSearchCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to search")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to search")
	cmd.Flags().IntVarP(&options.topK, "top-k", "k", defaultSearchTopK, "number of results to return")
	cmd.Flags().Var(&options.inputs, "inputs", "query inputs for search (inline JSON, ./path.json, or '-' for stdin)")
	cmd.Flags().StringVar(&options.id, "id", "", "use an existing record's vector by ID for the query")
	cmd.Flags().Var(&options.fields, "fields", "fields to return in results (JSON string array, ./path.json, or '-' for stdin)")
	cmd.Flags().StringVar(&options.body, "body", "", "request body JSON (inline, ./path.json, or '-' for stdin; only one argument may use stdin)")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")

	_ = cmd.MarkFlagRequired("index-name")

	return cmd
}

func runSearchCmd(ctx context.Context, options searchCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)

	var body *pinecone.SearchRecordsRequest
	if options.body != "" {
		b, src, err := argio.DecodeJSONArg[pinecone.SearchRecordsRequest](options.body)
		if err != nil {
			msg.FailMsg("Failed to parse search body (%s): %s", style.Emphasis(src.Label), err)
			exit.Errorf(err, "Failed to parse search body (%s)", src.Label)
		} else if b != nil {
			body = b
			if options.id == "" && b.Query.Id != nil {
				options.id = *b.Query.Id
			}
			if options.inputs == nil && b.Query.Inputs != nil {
				options.inputs = flags.JSONObject(*b.Query.Inputs)
			}
			if options.topK == defaultSearchTopK && b.Query.TopK > 0 {
				options.topK = int(b.Query.TopK)
			}
			if len(options.fields) == 0 && b.Fields != nil {
				options.fields = append(options.fields, (*b.Fields)...)
			}
		}
	}

	if options.topK <= 0 {
		msg.FailMsg("Top-k must be greater than 0")
		exit.ErrorMsg("Invalid top-k value")
	}

	req := pinecone.SearchRecordsRequest{
		Query: pinecone.SearchRecordsQuery{
			TopK: int32(options.topK),
		},
	}

	if options.id != "" {
		req.Query.Id = &options.id
	}

	if options.inputs != nil {
		inputs := map[string]interface{}(options.inputs)
		req.Query.Inputs = &inputs
	}

	if len(options.fields) > 0 {
		fieldsCopy := make([]string, len(options.fields))
		copy(fieldsCopy, options.fields)
		req.Fields = &fieldsCopy
	}

	if body != nil {
		if req.Query.TopK == 0 && body.Query.TopK > 0 {
			req.Query.TopK = body.Query.TopK
		}
		if req.Query.Id == nil && body.Query.Id != nil {
			req.Query.Id = body.Query.Id
		}
		if req.Query.Inputs == nil && body.Query.Inputs != nil {
			req.Query.Inputs = body.Query.Inputs
		}
		if req.Query.Vector == nil && body.Query.Vector != nil {
			req.Query.Vector = body.Query.Vector
		}
		if req.Query.Filter == nil && body.Query.Filter != nil {
			req.Query.Filter = body.Query.Filter
		}
		if req.Query.MatchTerms == nil && body.Query.MatchTerms != nil {
			req.Query.MatchTerms = body.Query.MatchTerms
		}
		if req.Rerank == nil && body.Rerank != nil {
			req.Rerank = body.Rerank
		}
		if req.Fields == nil && body.Fields != nil {
			req.Fields = body.Fields
		}
	}

	if req.Query.Id == nil && req.Query.Inputs == nil && req.Query.Vector == nil && req.Query.MatchTerms == nil {
		msg.FailMsg("Provide --inputs, --id, or a body with vector/match_terms")
		exit.ErrorMsg("Missing query inputs for search")
	}

	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, options.namespace)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	resp, err := ic.SearchRecords(ctx, &req)
	if err != nil {
		exit.Error(err, "Failed to search records")
	}

	if options.json {
		json := text.IndentJSON(resp)
		pcio.Println(json)
	} else {
		presenters.PrintSearchRecordsTable(resp)
	}
}
