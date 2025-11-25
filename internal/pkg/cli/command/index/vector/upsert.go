package vector

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/pinecone-io/cli/internal/pkg/utils/argio"
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
		Use:   "upsert",
		Short: "Upsert vectors into an index from a JSON/JSONL file",
		Long: help.Long(`
			Upsert vectors into an index namespace from a JSON or JSONL payload.
			
			The request --body may be a JSON object containing "vectors": [...] or a JSONL stream of Vector objects.
			Control batch size with --batch-size. Bodies can be inline JSON, loaded via @file, or read from stdin with @-.
		`),
		Example: help.Examples(`
			pc index vector upsert --index-name my-index --namespace my-namespace --body @./vectors.json
			pc index vector upsert --index-name my-index --namespace my-namespace --body @./vectors.jsonl
			cat payload.json | pc index vector upsert --index-name my-index --namespace my-namespace --body @-
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runUpsertCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of index to upsert into")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to upsert into")
	cmd.Flags().StringVar(&options.body, "body", "", "request body JSON or JSONL (inline, @path.json[l], or @- for stdin; only one argument may use stdin; max size: see PC_CLI_MAX_JSON_BYTES)")
	cmd.Flags().IntVarP(&options.batchSize, "batch-size", "b", 500, "size of batches to upsert (default: 500)")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	_ = cmd.MarkFlagRequired("index-name")
	_ = cmd.MarkFlagRequired("body")

	return cmd
}

func runUpsertCmd(ctx context.Context, options upsertCmdOptions) {
	b, src, err := argio.ReadAll(options.body)
	if err != nil {
		msg.FailMsg("Failed to read upsert body (%s): %s", style.Emphasis(src.Label), err)
		exit.Error(err, "Failed to read upsert body")
	}

	payload, err := parseUpsertBody(b)
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

func parseUpsertBody(b []byte) (*upsertBody, error) {
	var payload upsertBody
	// First try and decode as JSON: {"vectors":[...]}
	{
		dec := json.NewDecoder(bytes.NewReader(b))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&payload); err == nil && len(payload.Vectors) > 0 {
			return &payload, nil
		}
	}

	// Fallback: JSONL/stream of pinecone.Vector values
	var vectors []pinecone.Vector
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	for {
		var v pinecone.Vector
		if err := dec.Decode(&v); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		vectors = append(vectors, v)
	}
	if len(vectors) == 0 {
		return nil, io.EOF
	}
	return &upsertBody{Vectors: vectors}, nil
}
