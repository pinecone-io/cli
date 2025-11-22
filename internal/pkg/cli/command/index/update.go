package index

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/bodyutil"
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

type updateBody struct {
	Id           string                 `json:"id"`
	Values       []float32              `json:"values"`
	SparseValues *pinecone.SparseValues `json:"sparse_values"`
	Metadata     map[string]any         `json:"metadata"`
	Filter       map[string]any         `json:"filter"`
	DryRun       *bool                  `json:"dry_run"`
}

type updateCmdOptions struct {
	indexName     string
	namespace     string
	id            string
	values        flags.Float32List
	sparseIndices flags.UInt32List
	sparseValues  flags.Float32List
	metadata      flags.JSONObject
	filter        flags.JSONObject
	dryRun        bool
	body          string
	json          bool
}

func NewUpdateCmd() *cobra.Command {
	options := updateCmdOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a vector by ID, or a set of vectors by metadata filter",
		Example: help.Examples(`
			pc index update --index-name my-index --id doc-123 --values '[0.1, 0.2, 0.3]'
			pc index update --index-name my-index --id doc-123 --sparse-indices @./indices.json --sparse-values @./values.json
			pc index update --index-name my-index --id doc-123 --metadata '{"genre":"sci-fi"}'
			pc index update --index-name my-index --filter '{"genre":"sci-fi"}' --metadata '{"genre":"fantasy"}' --dry-run
			pc index update --index-name my-index --body @./update.json
			cat update.json | pc index update --index-name my-index --body @-
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runUpdateCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of the index to update")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to update the vector in")
	cmd.Flags().StringVar(&options.id, "id", "", "ID of the vector to update")
	cmd.Flags().Var(&options.values, "values", "values to update the vector with (inline JSON array, @path.json, or @- for stdin)")
	cmd.Flags().Var(&options.sparseIndices, "sparse-indices", "sparse indices to update the vector with (inline JSON array, @path.json, or @- for stdin)")
	cmd.Flags().Var(&options.sparseValues, "sparse-values", "sparse values to update the vector with (inline JSON array, @path.json, or @- for stdin)")
	cmd.Flags().Var(&options.metadata, "metadata", "metadata to update the vector with (inline JSON, @path.json, or @- for stdin)")
	cmd.Flags().Var(&options.filter, "filter", "filter to update the vectors with (inline JSON, @path.json, or @- for stdin)")
	cmd.Flags().BoolVar(&options.dryRun, "dry-run", false, "do not update the vectors, just return the number of vectors that would be updated")
	cmd.Flags().StringVar(&options.body, "body", "", "request body JSON (inline, @path.json, or @- for stdin; only one argument may use stdin)")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	_ = cmd.MarkFlagRequired("index-name")
	cmd.MarkFlagsMutuallyExclusive("id", "filter")

	return cmd
}

func runUpdateCmd(ctx context.Context, options updateCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)

	// Apply body overlay if provided
	if options.body != "" {
		if b, _, err := bodyutil.DecodeBodyArgs[updateBody](options.body); err != nil {
			exit.Error(err, "Failed to parse update body")
		} else if b != nil {
			if options.id == "" && b.Id != "" {
				options.id = b.Id
			}
			if len(options.values) == 0 && len(b.Values) > 0 {
				options.values = b.Values
			}
			if (len(options.sparseIndices) == 0 && len(options.sparseValues) == 0) && b.SparseValues != nil {
				options.sparseIndices = b.SparseValues.Indices
				options.sparseValues = b.SparseValues.Values
			}
			if options.filter == nil && b.Filter != nil {
				options.filter = b.Filter
			}
			if options.metadata == nil && b.Metadata != nil {
				options.metadata = b.Metadata
			}
			if b.DryRun != nil {
				options.dryRun = *b.DryRun
			}
		}
	}

	// Validate update by ID or metadata filter
	if options.id == "" && options.filter == nil {
		msg.FailMsg("Either --id or --filter must be provided")
		exit.ErrorMsg("Either --id or --filter must be provided")
	}

	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, options.namespace)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	// Update vector by ID
	if options.id != "" {
		metadata, err := pinecone.NewMetadata(options.metadata)
		if err != nil {
			msg.FailMsg("Failed to create metadata: %s", err)
			exit.Errorf(err, "Failed to create metadata")
		}

		var sparseValues *pinecone.SparseValues
		if len(options.sparseIndices) > 0 || len(options.sparseValues) > 0 {
			sparseValues = &pinecone.SparseValues{
				Indices: options.sparseIndices,
				Values:  options.sparseValues,
			}
		}

		err = ic.UpdateVector(ctx, &pinecone.UpdateVectorRequest{
			Id:           options.id,
			Values:       options.values,
			SparseValues: sparseValues,
			Metadata:     metadata,
		})
		if err != nil {
			msg.FailMsg("Failed to update vector ID: %s - %v", options.id, err)
			exit.Errorf(err, "Failed to update vector ID: %s", options.id)
		}

		if !options.json {
			msg.SuccessMsg("Vector ID: %s updated successfully", options.id)
		}
		return
	}

	// Update vectors by metadata filter
	if options.filter != nil {
		filter, err := pinecone.NewMetadataFilter(options.filter)
		if err != nil {
			msg.FailMsg("Failed to create filter: %s", err)
			exit.Errorf(err, "Failed to create filter")
		}

		metadata, err := pinecone.NewMetadata(options.metadata)
		if err != nil {
			msg.FailMsg("Failed to create metadata: %s", err)
			exit.Errorf(err, "Failed to create metadata")
		}

		var dryRun *bool
		if options.dryRun {
			dryRun = &options.dryRun
		}

		resp, err := ic.UpdateVectorsByMetadata(ctx, &pinecone.UpdateVectorsByMetadataRequest{
			Filter:   filter,
			Metadata: metadata,
			DryRun:   dryRun,
		})
		if err != nil {
			msg.FailMsg("Failed to update vectors by metadata: %s - %v", filter.String(), err)
			exit.Errorf(err, "Failed to update vectors by metadata: %s", filter.String())
		}

		if !options.json {
			msg.SuccessMsg("Updated %d vectors by metadata filter: %s", resp.MatchedRecords, filter.String())
			presenters.PrintUpdateVectorsByMetadataTable(resp)
		} else {
			json := text.IndentJSON(resp)
			pcio.Println(json)
		}
		return
	}
}
