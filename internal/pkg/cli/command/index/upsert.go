package index

import (
	"context"
	"encoding/json"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

type upsertCmdOptions struct {
	file      string
	name      string
	namespace string
	batchSize int
	json      bool
}

type upsertFile struct {
	Vectors   []upsertVector `json:"vectors"`
	Namespace string         `json:"namespace"`
}

type upsertVector struct {
	ID           string         `json:"id"`
	Values       []float32      `json:"values"`
	SparseValues *sparseValues  `json:"sparse_values,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

type sparseValues struct {
	Indices []uint32  `json:"indices"`
	Values  []float32 `json:"values"`
}

func NewUpsertCmd() *cobra.Command {
	options := upsertCmdOptions{}

	cmd := &cobra.Command{
		Use:   "upsert [file]",
		Short: "Upsert vectors into an index from a JSON file",
		Example: help.Examples(`
			pc index upsert --name my-index --namespace my-namespace ./vectors.json
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runUpsertCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of index to upsert into")
	cmd.Flags().StringVar(&options.namespace, "namespace", "", "namespace to upsert into")
	cmd.Flags().StringVarP(&options.file, "file", "f", "", "file to upsert from")
	cmd.Flags().IntVarP(&options.batchSize, "batch-size", "b", 1000, "size of batches to upsert (default: 1000)")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runUpsertCmd(ctx context.Context, options upsertCmdOptions) {
	filePath := options.file
	raw, err := os.ReadFile(filePath)
	if err != nil {
		msg.FailMsg("Failed to read file %s: %s", style.Emphasis(filePath), err)
		exit.Errorf(err, "Failed to read file %s", filePath)
	}

	var payload upsertFile
	if err := json.Unmarshal(raw, &payload); err != nil {
		msg.FailMsg("Failed to parse JSON from %s: %s", style.Emphasis(filePath), err)
		exit.Error(err, "Failed to parse JSON for upsert")
	}

	// Default namespace
	ns := payload.Namespace
	if options.namespace != "" {
		ns = options.namespace
	}
	if ns == "" {
		ns = "__default__"
	}
	// Get IndexConnection
	pc := sdk.NewPineconeClient(ctx)
	ic, err := sdk.NewIndexConnection(ctx, pc, options.name, ns)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	if len(payload.Vectors) == 0 {
		msg.FailMsg("No vectors found in %s", style.Emphasis(filePath))
		exit.ErrorMsg("No vectors provided for upsert")
	}

	// Map to SDK types
	mapped := make([]*pinecone.Vector, 0, len(payload.Vectors))
	for _, v := range payload.Vectors {
		values := v.Values
		metadata, err := pinecone.NewMetadata(v.Metadata)
		if err != nil {
			msg.FailMsg("Failed to parse metadata: %s", err)
			exit.Error(err, "Failed to parse metadata")
		}

		var vector pinecone.Vector
		vector.Id = v.ID
		if v.Values != nil {
			vector.Values = &values
		}
		if v.SparseValues != nil {
			vector.SparseValues = &pinecone.SparseValues{
				Indices: v.SparseValues.Indices,
				Values:  v.SparseValues.Values,
			}
		}
		vector.Metadata = metadata
		mapped = append(mapped, &vector)
	}

	batches := make([][]*pinecone.Vector, 0, (len(mapped)+options.batchSize-1)/options.batchSize)
	for i := 0; i < len(mapped); i += options.batchSize {
		end := i + options.batchSize
		if end > len(mapped) {
			end = len(mapped)
		}
		batches = append(batches, mapped[i:end])
	}

	for i, batch := range batches {
		resp, err := ic.UpsertVectors(ctx, batch)
		if err != nil {
			msg.FailMsg("Failed to upsert %d vectors in batch %d: %s", len(batch), i+1, err)
			exit.Errorf(err, "Failed to upsert %d vectors in batch %d", len(batch), i+1)
		} else {
			if options.json {
				json := text.IndentJSON(resp)
				pcio.Println(json)
			} else {
				msg.SuccessMsg("Upserted %d vectors into namespace %s in %d batches", len(batch), ns, i+1)
			}
		}
	}
}
