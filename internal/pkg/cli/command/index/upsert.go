package index

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/utils/bodyutil"
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

type upsertBody struct {
	Vectors []pinecone.Vector `json:"vectors"`
}

type upsertCmdOptions struct {
	body      string
	indexName string
	namespace string
	batchSize int
	json      bool
}

func NewUpsertCmd() *cobra.Command {
	options := upsertCmdOptions{}

	cmd := &cobra.Command{
		Use:   "upsert [file]",
		Short: "Upsert vectors into an index from a JSON file",
		Example: help.Examples(`
			pc index upsert --index-name my-index --namespace my-namespace ./vectors.json
			pc index upsert --index-name my-index --namespace my-namespace --file - < ./vectors.json
			pc index upsert --index-name my-index --body @./payload.json
			cat payload.json | pc index upsert --index-name my-index --body @-
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runUpsertCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of index to upsert into")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to upsert into")
	cmd.Flags().StringVar(&options.body, "body", "", "request body JSON (inline, @path.json, or @- for stdin; only one argument may use stdin; max size: see PC_CLI_MAX_JSON_BYTES)")
	cmd.Flags().IntVarP(&options.batchSize, "batch-size", "b", 1000, "size of batches to upsert (default: 1000)")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	_ = cmd.MarkFlagRequired("index-name")
	_ = cmd.MarkFlagRequired("body")

	return cmd
}

func runUpsertCmd(ctx context.Context, options upsertCmdOptions) {
	var payload *upsertBody
	payload, src, err := bodyutil.DecodeBodyArgs[upsertBody](options.body)
	if err != nil {
		msg.FailMsg("Failed to parse upsert body (%s): %s", style.Emphasis(src.Label), err)
		exit.Error(err, "Failed to parse upsert body")
	}

	// Get IndexConnection
	pc := sdk.NewPineconeClient(ctx)
	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, options.namespace)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	if len(payload.Vectors) == 0 {
		msg.FailMsg("No vectors found in %s", style.Emphasis(src.Label))
		exit.ErrorMsg("No vectors provided for upsert")
	}

	// Map to SDK types
	mapped := make([]*pinecone.Vector, 0, len(payload.Vectors))
	for _, v := range payload.Vectors {
		values := v.Values

		var vector pinecone.Vector
		vector.Id = v.Id
		vector.Metadata = v.Metadata

		if values != nil {
			vector.Values = values
		}

		if v.SparseValues != nil {
			vector.SparseValues = &pinecone.SparseValues{
				Indices: v.SparseValues.Indices,
				Values:  v.SparseValues.Values,
			}
		}

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
				msg.SuccessMsg("Upserted %d vectors into namespace %s (batch %d of %d)", len(batch), options.namespace, i+1, len(batches))
			}
		}
	}
}
