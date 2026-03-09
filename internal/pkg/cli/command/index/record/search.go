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
	indexName     string
	namespace     string
	topK          int
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
			runSearchCmd(cmd, options)
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
	if options.matchTerms != nil {
		b, err := json.Marshal(options.matchTerms)
		if err != nil {
			msg.FailMsg("Failed to encode --match-terms value: %s", err)
			exit.Error(err, "Failed to encode --match-terms value")
		}
		var matchTerms pinecone.SearchMatchTerms
		if err := json.Unmarshal(b, &matchTerms); err != nil {
			msg.FailMsg("Failed to parse --match-terms value: %s", err)
			exit.Error(err, "Failed to parse --match-terms value")
		}
		req.Query.MatchTerms = &matchTerms
	}
	if len(options.vector) > 0 || len(options.sparseIndices) > 0 {
		if len(options.sparseIndices) != len(options.sparseValues) {
			msg.FailMsg("--sparse-indices and --sparse-values must be the same length")
			exit.ErrorMsg("--sparse-indices and --sparse-values must be the same length")
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

	// Merge --body into req. For all fields that have a dedicated flag, the flag
	// takes precedence; body only fills in values that weren't explicitly set.
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
			vectorFlagsSet := cmd.Flags().Changed("vector") || cmd.Flags().Changed("sparse-indices") || cmd.Flags().Changed("sparse-values")
			if !vectorFlagsSet && b.Query.Vector != nil {
				req.Query.Vector = b.Query.Vector
			}
			if !cmd.Flags().Changed("match-terms") && b.Query.MatchTerms != nil {
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

	if req.Query.Id == nil && req.Query.Inputs == nil && req.Query.Vector == nil {
		msg.FailMsg("Provide a query via --inputs, --id, --vector, --sparse-indices/--sparse-values, or a --body")
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
