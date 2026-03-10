package record

import (
	"context"
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
	indexName     string
	namespace     string
	topK          int
	topKExplicit  bool // true when --top-k was explicitly passed (vs. the default)
	inputs        flags.JSONObject
	filter        flags.JSONObject
	rerank        flags.JSONObject
	id            string
	vector        flags.Float32List
	sparseIndices flags.Int32List
	sparseValues  flags.Float32List
	matchTerms    flags.JSONObject
	fields        flags.StringList
	body          string
	json          bool
}

func NewSearchCmd() *cobra.Command {
	options := searchCmdOptions{
		topK: defaultSearchTopK,
	}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search records in an index",
		Long: help.Long(`
			Search an index and return the most similar records along with their similarity
			scores. Exactly one primary query mode must be provided:

			  --inputs   query text (requires an index with integrated embedding)
			  --id       use an existing record's vector as the query
			  --vector   provide a dense vector directly (all index types)
			  --sparse-indices + --sparse-values   sparse vector (all index types)

			Combine --vector with --sparse-indices/--sparse-values for hybrid search.

			Optional modifiers:
			  --filter      metadata filter applied before ranking
			  --rerank      re-score results using an inference model
			  --fields      restrict which fields are returned
			  --top-k       number of results to return (default: 10)
			  --namespace   target a specific namespace
			  --match-terms keyword terms that must appear in results; only supported for
			                sparse indexes with integrated embedding using the
			                pinecone-sparse-english-v0 model; requires --inputs

			Use --body to pass a full request object. Flags take precedence over --body
			when both specify the same field.
		`),
		Example: help.Examples(`
			# Text search (integrated embedding indexes only)
			pc index record search --index-name my-index --inputs '{"text":"disease prevention"}'
			pc index record search --index-name my-index --inputs '{"text":"disease prevention"}' --namespace my-namespace --top-k 5
			pc index record search --index-name my-index --inputs '{"text":"disease prevention"}' --filter '{"category":"health"}'
			pc index record search --index-name my-index --inputs '{"text":"disease prevention"}' --fields '["_id","chunk_text","category"]'
			pc index record search --index-name my-index --inputs '{"text":"disease prevention"}' --rerank '{"model":"bge-reranker-v2-m3","rank_fields":["chunk_text"]}'
			pc index record search --index-name my-index --inputs '{"text":"disease prevention"}' --match-terms '{"terms":["vaccine","prevention"]}'
			echo '{"text":"disease prevention"}' | pc index record search --index-name my-index --inputs -

			# Search by record ID (all index types)
			pc index record search --index-name my-index --id rec-123
			pc index record search --index-name my-index --id rec-123 --fields '["_id","chunk_text"]'

			# Vector search (all index types)
			pc index record search --index-name my-index --vector '[0.1,0.2,0.3]'
			pc index record search --index-name my-index --sparse-indices '[1,5,9]' --sparse-values '[0.4,0.8,0.2]'
			pc index record search --index-name my-index --vector '[0.1,0.2,0.3]' --sparse-indices '[1,5,9]' --sparse-values '[0.4,0.8,0.2]'

			# Full request body
			pc index record search --index-name my-index --body ./search.json
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			options.topKExplicit = cmd.Flags().Changed("top-k")
			pc := sdk.NewPineconeClient(ctx)
			ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, options.namespace)
			if err != nil {
				msg.FailMsg("Failed to create index connection: %s", err)
				exit.Error(err, "Failed to create index connection")
			}
			if err := runSearchCmd(ctx, ic, options); err != nil {
				msg.FailMsg("%s", err)
				exit.Error(err, "search failed")
			}
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to search")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to search")
	cmd.Flags().IntVarP(&options.topK, "top-k", "k", defaultSearchTopK, "number of results to return")
	cmd.Flags().Var(&options.inputs, "inputs", "query inputs for search (inline JSON, ./path.json, or '-' for stdin); requires integrated embedding")
	cmd.Flags().Var(&options.filter, "filter", "metadata filter (inline JSON, ./path.json, or '-' for stdin)")
	cmd.Flags().Var(&options.rerank, "rerank", "rerank results (inline JSON, ./path.json, or '-' for stdin); required fields: model (string), rank_fields (string array)")
	cmd.Flags().StringVar(&options.id, "id", "", "use an existing record's vector by ID for the query")
	cmd.Flags().VarP(&options.vector, "vector", "v", "dense vector values to search against (inline JSON float32 array, ./path.json, or '-' for stdin)")
	cmd.Flags().Var(&options.sparseIndices, "sparse-indices", "sparse vector indices (inline JSON int32 array, ./path.json, or '-' for stdin)")
	cmd.Flags().Var(&options.sparseValues, "sparse-values", "sparse vector values (inline JSON float32 array, ./path.json, or '-' for stdin)")
	cmd.Flags().Var(&options.matchTerms, "match-terms", "keyword terms filter for sparse integrated indexes (inline JSON, ./path.json, or '-' for stdin); required field: terms (string array); optional: strategy (default: \"all\"); requires --inputs")
	cmd.Flags().Var(&options.fields, "fields", "fields to return in results (inline JSON string array, ./path.json, or '-' for stdin)")
	cmd.Flags().StringVar(&options.body, "body", "", "request body JSON (inline, ./path.json, or '-' for stdin; only one argument may use stdin)")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")

	_ = cmd.MarkFlagRequired("index-name")
	cmd.MarkFlagsMutuallyExclusive("inputs", "id", "vector")
	cmd.MarkFlagsMutuallyExclusive("inputs", "id", "sparse-values")
	cmd.MarkFlagsRequiredTogether("sparse-indices", "sparse-values")

	return cmd
}

func runSearchCmd(ctx context.Context, ic RecordService, options searchCmdOptions) error {
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
			return pcio.Errorf("failed to encode --rerank value: %w", err)
		}
		var rerank pinecone.SearchRecordsRerank
		if err := json.Unmarshal(b, &rerank); err != nil {
			return pcio.Errorf("failed to parse --rerank value: %w", err)
		}
		req.Rerank = &rerank
	}
	if options.matchTerms != nil {
		b, err := json.Marshal(options.matchTerms)
		if err != nil {
			return pcio.Errorf("failed to encode --match-terms value: %w", err)
		}
		var matchTerms pinecone.SearchMatchTerms
		if err := json.Unmarshal(b, &matchTerms); err != nil {
			return pcio.Errorf("failed to parse --match-terms value: %w", err)
		}
		req.Query.MatchTerms = &matchTerms
	}
	if len(options.vector) > 0 || len(options.sparseIndices) > 0 {
		if len(options.sparseIndices) != len(options.sparseValues) {
			return pcio.Errorf("--sparse-indices and --sparse-values must be the same length")
		}
		sv := &pinecone.SearchRecordsVector{}
		if len(options.vector) > 0 {
			values := []float32(options.vector)
			sv.Values = &values
		}
		if len(options.sparseIndices) > 0 {
			indices := []int32(options.sparseIndices)
			sparseVals := []float32(options.sparseValues)
			sv.SparseIndices = &indices
			sv.SparseValues = &sparseVals
		}
		req.Query.Vector = sv
	}

	// Merge --body into req. Flags take precedence: a flag's value in req is
	// already non-nil/non-zero if the flag was set, so we only overwrite when
	// the field is still unset. top-k is the exception — its default is non-zero,
	// so we track whether it was explicitly passed via options.topKExplicit.
	if options.body != "" {
		b, src, err := argio.DecodeJSONArg[pinecone.SearchRecordsRequest](options.body)
		if err != nil {
			return pcio.Errorf("failed to parse search body (%s): %w", style.Emphasis(src.Label), err)
		}
		if b != nil {
			if !options.topKExplicit && b.Query.TopK > 0 {
				req.Query.TopK = b.Query.TopK
			}
			if req.Query.Id == nil && b.Query.Id != nil {
				req.Query.Id = b.Query.Id
			}
			if req.Query.Inputs == nil && b.Query.Inputs != nil {
				req.Query.Inputs = b.Query.Inputs
			}
			if req.Query.Filter == nil && b.Query.Filter != nil {
				req.Query.Filter = b.Query.Filter
			}
			if req.Fields == nil && b.Fields != nil {
				req.Fields = b.Fields
			}
			if req.Query.Vector == nil && b.Query.Vector != nil {
				req.Query.Vector = b.Query.Vector
			}
			if req.Query.MatchTerms == nil && b.Query.MatchTerms != nil {
				req.Query.MatchTerms = b.Query.MatchTerms
			}
			if req.Rerank == nil && b.Rerank != nil {
				req.Rerank = b.Rerank
			}
		}
	}

	if req.Query.TopK <= 0 {
		return pcio.Errorf("top-k must be greater than 0")
	}

	if req.Query.Id == nil && req.Query.Inputs == nil && req.Query.Vector == nil {
		return pcio.Errorf("provide a query via --inputs, --id, --vector, --sparse-indices/--sparse-values, or a --body")
	}

	resp, err := ic.SearchRecords(ctx, &req)
	if err != nil {
		return pcio.Errorf("failed to search records: %w", err)
	}

	if options.json {
		pcio.Println(text.IndentJSON(resp))
	} else {
		presenters.PrintSearchRecordsTable(resp)
	}

	return nil
}
