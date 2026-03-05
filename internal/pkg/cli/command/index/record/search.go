package record

import (
	"encoding/json"

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
	filter    flags.JSONObject
	rerank    flags.JSONObject
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
			Narrow results with --filter (metadata filter as a JSON object).
			Rerank results with --rerank (JSON object with required fields: model, rank_fields).
			Use --body to supply a full request body for advanced parameters like
			vector overrides or match_terms.

			When a flag and --body both specify the same field, the flag takes precedence.
		`),
		Example: help.Examples(`
			pc index record search --index-name my-index --namespace my-namespace --inputs '{"text":"find similar"}'
			pc index record search --index-name my-index --inputs '{"text":"disease prevention"}' --filter '{"category":"health"}'
			pc index record search --index-name my-index --inputs '{"text":"find similar"}' --rerank '{"model":"bge-reranker-v2-m3","rank_fields":["chunk_text"]}'
			echo '{"text":"disease prevention"}' | pc index record search --index-name my-index --inputs -
			pc index record search --index-name my-index --namespace my-namespace --body ./search.json
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runSearchCmd(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to search")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to search")
	cmd.Flags().IntVarP(&options.topK, "top-k", "k", defaultSearchTopK, "number of results to return")
	cmd.Flags().Var(&options.inputs, "inputs", "query inputs for search (inline JSON, ./path.json, or '-' for stdin)")
	cmd.Flags().Var(&options.filter, "filter", "metadata filter (inline JSON, ./path.json, or '-' for stdin)")
	cmd.Flags().Var(&options.rerank, "rerank", "rerank results (inline JSON, ./path.json, or '-' for stdin); required fields: model (string), rank_fields (string array)")
	cmd.Flags().StringVar(&options.id, "id", "", "use an existing record's vector by ID for the query")
	cmd.Flags().Var(&options.fields, "fields", "fields to return in results (JSON string array, ./path.json, or '-' for stdin)")
	cmd.Flags().StringVar(&options.body, "body", "", "request body JSON (inline, ./path.json, or '-' for stdin; only one argument may use stdin)")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")

	_ = cmd.MarkFlagRequired("index-name")

	return cmd
}

func runSearchCmd(cmd *cobra.Command, options searchCmdOptions) {
	ctx := cmd.Context()
	pc := sdk.NewPineconeClient(ctx)

	// Build req from flags.
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
	if options.filter != nil {
		filter := map[string]interface{}(options.filter)
		req.Query.Filter = &filter
	}
	if len(options.fields) > 0 {
		fieldsCopy := make([]string, len(options.fields))
		copy(fieldsCopy, options.fields)
		req.Fields = &fieldsCopy
	}
	if options.rerank != nil {
		b, err := json.Marshal(options.rerank)
		if err != nil {
			msg.FailMsg("Failed to encode --rerank value: %s", err)
			exit.Error(err, "Failed to encode --rerank value")
		}
		var rerank pinecone.SearchRecordsRerank
		if err := json.Unmarshal(b, &rerank); err != nil {
			msg.FailMsg("Failed to parse --rerank value: %s", err)
			exit.Error(err, "Failed to parse --rerank value")
		}
		req.Rerank = &rerank
	}

	// Merge --body into req. For fields that have a dedicated flag, the flag
	// takes precedence; body only fills in values that weren't explicitly set.
	// Fields with no dedicated flag (Vector, MatchTerms, Rerank) always come
	// from body.
	if options.body != "" {
		b, src, err := argio.DecodeJSONArg[pinecone.SearchRecordsRequest](options.body)
		if err != nil {
			msg.FailMsg("Failed to parse search body (%s): %s", style.Emphasis(src.Label), err)
			exit.Errorf(err, "Failed to parse search body (%s)", src.Label)
		}
		if b != nil {
			if !cmd.Flags().Changed("top-k") && b.Query.TopK > 0 {
				req.Query.TopK = b.Query.TopK
			}
			if !cmd.Flags().Changed("id") && b.Query.Id != nil {
				req.Query.Id = b.Query.Id
			}
			if !cmd.Flags().Changed("inputs") && b.Query.Inputs != nil {
				req.Query.Inputs = b.Query.Inputs
			}
			if !cmd.Flags().Changed("filter") && b.Query.Filter != nil {
				req.Query.Filter = b.Query.Filter
			}
			if !cmd.Flags().Changed("fields") && b.Fields != nil {
				req.Fields = b.Fields
			}
			if b.Query.Vector != nil {
				req.Query.Vector = b.Query.Vector
			}
			if b.Query.MatchTerms != nil {
				req.Query.MatchTerms = b.Query.MatchTerms
			}
			if !cmd.Flags().Changed("rerank") && b.Rerank != nil {
				req.Rerank = b.Rerank
			}
		}
	}

	if req.Query.TopK <= 0 {
		msg.FailMsg("Top-k must be greater than 0")
		exit.ErrorMsg("Invalid top-k value")
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
