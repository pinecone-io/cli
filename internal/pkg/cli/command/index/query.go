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

type queryCmdOptions struct {
	id              string
	vector          []float32
	sparseIndices   []int32
	sparseValues    []float32
	name            string
	namespace       string
	topK            uint32
	filter          string
	filterFile      string
	includeValues   bool
	includeMetadata bool
	json            bool
}

func NewQueryCmd() *cobra.Command {
	options := queryCmdOptions{}
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query an index by vector values",
		Example: help.Examples(`
		
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runQueryCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the index to query")
	cmd.Flags().StringVar(&options.namespace, "namespace", "", "index namespace to query")
	cmd.Flags().Uint32VarP(&options.topK, "top-k", "k", 10, "maximum number of results to return")
	cmd.Flags().StringVarP(&options.filter, "filter", "f", "", "metadata filter to apply to the query")
	cmd.Flags().StringVar(&options.filterFile, "filter-file", "", "file containing metadata filter to apply to the query")
	cmd.Flags().BoolVar(&options.includeValues, "include-values", false, "include vector values in the query results")
	cmd.Flags().BoolVar(&options.includeMetadata, "include-metadata", false, "include metadata in the query results")
	cmd.Flags().StringVarP(&options.id, "id", "i", "", "ID of the vector to query against")
	cmd.Flags().Float32SliceVarP(&options.vector, "vector", "v", []float32{}, "vector values to query against")
	cmd.Flags().Int32SliceVar(&options.sparseIndices, "sparse-indices", []int32{}, "sparse indices to query against")
	cmd.Flags().Float32SliceVar(&options.sparseValues, "sparse-values", []float32{}, "sparse values to query against")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	_ = cmd.MarkFlagRequired("name")
	cmd.MarkFlagsMutuallyExclusive("id", "vector", "sparse-values")

	return cmd
}

func runQueryCmd(ctx context.Context, options queryCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)

	// Default namespace
	ns := options.namespace
	if options.namespace != "" {
		ns = options.namespace
	}
	if ns == "" {
		ns = "__default__"
	}

	// Get IndexConnection
	ic, err := sdk.NewIndexConnection(ctx, pc, options.name, ns)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	var queryResponse *pinecone.QueryVectorsResponse

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
		sparseIndices, err := toUint32Slice(options.sparseIndices)
		if err != nil {
			exit.Error(err, "Failed to convert sparse indices to uint32")
		}

		req := &pinecone.QueryByVectorValuesRequest{
			Vector: options.vector,
			SparseValues: &pinecone.SparseValues{
				Indices: sparseIndices,
				Values:  options.sparseValues,
			},
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

func toUint32Slice(in []int32) ([]uint32, error) {
	out := make([]uint32, len(in))
	for i, v := range in {
		if v < 0 {
			return nil, pcio.Errorf("sparse indices must be non-negative")
		}
		out[i] = uint32(v)
	}
	return out, nil
}
