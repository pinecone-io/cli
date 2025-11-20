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

type updateCmdOptions struct {
	name          string
	namespace     string
	id            string
	values        []float32
	sparseIndices []int32
	sparseValues  []float32
	metadata      flags.JSONObject
	filter        flags.JSONObject
	dryRun        bool
	json          bool
}

func NewUpdateCmd() *cobra.Command {
	options := updateCmdOptions{}

	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update a vector by ID, or a set of vectors by metadata filter",
		Example: help.Examples(``),
		Run: func(cmd *cobra.Command, args []string) {
			runUpdateCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the index to update")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to update the vector in")
	cmd.Flags().StringVar(&options.id, "id", "", "ID of the vector to update")
	cmd.Flags().Float32SliceVar(&options.values, "values", []float32{}, "values to update the vector with")
	cmd.Flags().Int32SliceVar(&options.sparseIndices, "sparse-indices", []int32{}, "sparse indices to update the vector with")
	cmd.Flags().Float32SliceVar(&options.sparseValues, "sparse-values", []float32{}, "sparse values to update the vector with")
	cmd.Flags().Var(&options.metadata, "metadata", "metadata to update the vector with")
	cmd.Flags().Var(&options.filter, "filter", "filter to update the vectors with")
	cmd.Flags().BoolVar(&options.dryRun, "dry-run", false, "do not update the vectors, just return the number of vectors that would be updated")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runUpdateCmd(ctx context.Context, options updateCmdOptions) {
	pc := sdk.NewPineconeClient(ctx)

	// Default namespace
	ns := options.namespace
	if ns == "" {
		ns = "__default__"
	}

	// Validate update by ID or metadata filter
	if options.id == "" && options.filter == nil {
		msg.FailMsg("Either --id or --filter must be provided")
		exit.ErrorMsg("Either --id or --filter must be provided")
	}

	ic, err := sdk.NewIndexConnection(ctx, pc, options.name, ns)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	if options.id != "" && options.filter != nil {
		msg.FailMsg("ID and filter cannot be used together")
		exit.ErrorMsg("ID and filter cannot be used together")
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
			sparseIndices, err := toUint32Slice(options.sparseIndices)
			if err != nil {
				msg.FailMsg("Failed to convert sparse indices to uint32: %s", err)
				exit.Errorf(err, "Failed to convert sparse indices to uint32")
			}
			sparseValues = &pinecone.SparseValues{
				Indices: sparseIndices,
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
