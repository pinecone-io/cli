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

// QueryBody is the JSON payload schema for --body.
// Fields: id, vector, sparse_values (https://pkg.go.dev/github.com/pinecone-io/go-pinecone/v5/pinecone#SparseValues),
// filter, top_k, include_values, include_metadata.
type QueryBody struct {
	Id              string                 `json:"id"`
	Vector          []float32              `json:"vector"`
	SparseValues    *pinecone.SparseValues `json:"sparse_values"`
	Filter          map[string]any         `json:"filter"`
	TopK            *uint32                `json:"top_k"`
	IncludeValues   *bool                  `json:"include_values"`
	IncludeMetadata *bool                  `json:"include_metadata"`
}

type queryCmdOptions struct {
	id              string
	vector          flags.Float32List
	sparseIndices   flags.UInt32List
	sparseValues    flags.Float32List
	indexName       string
	namespace       string
	topK            uint32
	filter          flags.JSONObject
	includeValues   bool
	includeMetadata bool
	body            string
	json            bool
}

func NewQueryCmd() *cobra.Command {
	options := queryCmdOptions{}
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query an index by vector values",
		Long: help.Long(`
			Query vectors in an index by dense or sparse vector values, or by vector ID..

			Use --top-k to control result count and --include-values/--include-metadata to enrich results.
			An optional metadata filter can be used to filter the results.

			JSON inputs may be inline, loaded from ./file.json[l], or read from stdin with '-'.

			When providing sparse values, both --sparse-indices and --sparse-values must be present.
			A --body payload can pass id, vector, sparse_values, filter, top_k, include_values, and include_metadata.
		`),
		Example: help.Examples(`
			pc index vector query --index-name my-index --id doc-123 --top-k 10 --include-metadata
		
			pc index vector query --index-name my-index --vector '[0.1, 0.2, 0.3]' --top-k 25
			pc index vector query --index-name my-index --vector ./vector.json --top-k 25 --include-metadata
			jq -c '.embedding' doc.json | pc index vector query --index-name my-index --vector - --top-k 20
		
			pc index vector query --index-name my-index --sparse-indices ./indices.json --sparse-values ./values.json --top-k 15
		
			pc index vector query --index-name my-index --vector ./vector.json --filter '{"genre":{"$eq":"sci-fi"}}' --include-metadata
		
			pc index vector query --index-name my-index --body ./query.json
			cat query.json | pc index vector query --index-name my-index --body -
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runQueryCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to query")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "index namespace to query")
	cmd.Flags().Uint32VarP(&options.topK, "top-k", "k", 10, "maximum number of results to return")
	cmd.Flags().VarP(&options.filter, "filter", "f", "metadata filter to apply to the query (inline JSON, ./path.json, or '-' for stdin)")
	cmd.Flags().BoolVar(&options.includeValues, "include-values", false, "include vector values in the query results")
	cmd.Flags().BoolVar(&options.includeMetadata, "include-metadata", false, "include metadata in the query results")
	cmd.Flags().StringVarP(&options.id, "id", "i", "", "ID of the vector to query against")
	cmd.Flags().VarP(&options.vector, "vector", "v", "vector values to query against (inline JSON array, ./path.json, or '-' for stdin)")
	cmd.Flags().Var(&options.sparseIndices, "sparse-indices", "sparse indices to query against (inline JSON array, ./path.json, or '-' for stdin)")
	cmd.Flags().Var(&options.sparseValues, "sparse-values", "sparse values to query against (inline JSON array, ./path.json, or '-' for stdin)")
	cmd.Flags().StringVar(&options.body, "body", "", "request body JSON (inline, ./path.json, or '-' for stdin; only one argument may use stdin)")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	_ = cmd.MarkFlagRequired("index-name")
	cmd.MarkFlagsMutuallyExclusive("id", "vector", "sparse-values")

	return cmd
}

func runQueryCmd(ctx context.Context, options queryCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)

	// Apply body overlay if provided
	if options.body != "" {
		if b, src, err := argio.DecodeJSONArg[QueryBody](options.body); err != nil {
			msg.FailMsg("Failed to parse query body (%s): %s", style.Emphasis(src.Label), err)
			exit.Errorf(err, "Failed to parse query body (%s): %v", src.Label, err)
		} else if b != nil {
			if options.id == "" && b.Id != "" {
				options.id = b.Id
			}
			if len(options.vector) == 0 && len(b.Vector) > 0 {
				options.vector = b.Vector
			}
			if (len(options.sparseIndices) == 0 && len(options.sparseValues) == 0) && b.SparseValues != nil {
				options.sparseIndices = b.SparseValues.Indices
				options.sparseValues = b.SparseValues.Values
			}
			if options.filter == nil && b.Filter != nil {
				options.filter = b.Filter
			}
			if b.TopK != nil {
				options.topK = *b.TopK
			}
			if b.IncludeValues != nil {
				options.includeValues = *b.IncludeValues
			}
			if b.IncludeMetadata != nil {
				options.includeMetadata = *b.IncludeMetadata
			}
		}
	}

	if options.id == "" && options.vector == nil && options.sparseIndices == nil && options.sparseValues == nil && options.filter == nil {
		msg.FailMsg("Either --id, --vector, --sparse-indices & --sparse-values, or --filter must be provided")
		exit.ErrorMsg("Either --id, --vector, --sparse-indices & --sparse-values, or --filter must be provided")
	}

	// Get IndexConnection
	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, options.namespace)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	var queryResponse *pinecone.QueryVectorsResponse

	// Build metadata filter if provided
	var filter *pinecone.MetadataFilter
	if options.filter != nil {
		filter, err = pinecone.NewMetadataFilter(options.filter)
		if err != nil {
			msg.FailMsg("Failed to create filter: %s", err)
			exit.Errorf(err, "Failed to create filter")
		}
	}

	// Query by vector ID
	if options.id != "" {
		req := &pinecone.QueryByVectorIdRequest{
			VectorId:        options.id,
			TopK:            options.topK,
			IncludeValues:   options.includeValues,
			IncludeMetadata: options.includeMetadata,
			MetadataFilter:  filter,
		}

		queryResponse, err = ic.QueryByVectorId(ctx, req)
		if err != nil {
			exit.Error(err, "Failed to query by vector ID")
		}
	}

	// Query by vector values
	if len(options.vector) > 0 || len(options.sparseIndices) > 0 || len(options.sparseValues) > 0 {
		var sparse *pinecone.SparseValues

		// Only include sparse values if the user provided them
		if len(options.sparseIndices) > 0 || len(options.sparseValues) > 0 {
			if len(options.sparseIndices) == 0 || len(options.sparseValues) == 0 {
				exit.Errorf(nil, "both --sparse-indices and --sparse-values are required when specifying sparse values")
			}
			if len(options.sparseIndices) != len(options.sparseValues) {
				exit.Errorf(nil, "--sparse-indices and --sparse-values must be the same length")
			}
			sparse = &pinecone.SparseValues{
				Indices: options.sparseIndices,
				Values:  options.sparseValues,
			}
		}

		req := &pinecone.QueryByVectorValuesRequest{
			Vector:          options.vector,
			SparseValues:    sparse,
			TopK:            options.topK,
			IncludeValues:   options.includeValues,
			IncludeMetadata: options.includeMetadata,
			MetadataFilter:  filter,
		}

		queryResponse, err = ic.QueryByVectorValues(ctx, req)
		if err != nil {
			exit.Error(err, "Failed to query by vector values")
		}
	}

	if options.json {
		json := text.IndentJSON(queryResponse)
		pcio.Println(json)
	} else {
		presenters.PrintQueryVectorsTable(queryResponse)
	}
}
