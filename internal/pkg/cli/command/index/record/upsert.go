package record

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

// UpsertRecordsBody is the JSON payload for --body/--file.
// It accepts either {"records": [...]} where each element is an IntegratedRecord,
// a JSON array of IntegratedRecord objects, or a JSONL stream of IntegratedRecord objects.
type UpsertRecordsBody struct {
	Records []pinecone.IntegratedRecord `json:"records"`
}

type upsertCmdOptions struct {
	file      string
	indexName string
	namespace string
	batchSize int
	json      bool
}

func NewUpsertCmd() *cobra.Command {
	options := upsertCmdOptions{}

	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "Upsert records into an index from a JSON/JSONL payload",
		Long: help.Long(`
			Upsert records into an index namespace from a JSON or JSONL payload.

			The request --body/--file may be a JSON object containing "records": [...]
			(a list of IntegratedRecord objects), a raw JSON array of records, or a
			JSONL stream of IntegratedRecord objects. Bodies can be inline JSON,
			loaded from ./file.json[l], or read from stdin with '-'.

			Body schema: UpsertRecordsBody (records shaped like pinecone.IntegratedRecord:
			https://pkg.go.dev/github.com/pinecone-io/go-pinecone/v5/pinecone#IntegratedRecord)
		`),
		Example: help.Examples(`
			pc index record upsert --index-name my-index --namespace my-namespace --body ./records.json
			pc index record upsert --index-name my-index --namespace my-namespace --body ./records.jsonl
			cat records.jsonl | pc index record upsert --index-name my-index --namespace my-namespace --body -
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runUpsertCmd(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of index to upsert into")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to upsert into")
	cmd.Flags().StringVar(&options.file, "file", "", "request body JSON or JSONL (inline, ./path.json[l], or '-' for stdin; only one argument may use stdin)")
	cmd.Flags().StringVar(&options.file, "body", "", "alias for --file")
	cmd.Flags().IntVarP(&options.batchSize, "batch-size", "b", 500, "size of batches to upsert (default: 500)")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")

	_ = cmd.MarkFlagRequired("index-name")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runUpsertCmd(ctx context.Context, options upsertCmdOptions) {
	b, src, err := argio.ReadAll(options.file)
	if err != nil {
		msg.FailMsg("Failed to read upsert body (%s): %s", style.Emphasis(src.Label), err)
		exit.Error(err, "Failed to read upsert body")
	}

	payload, err := parseUpsertRecordsBody(b)
	if err != nil {
		msg.FailMsg("Failed to parse upsert body (%s): %s", style.Emphasis(src.Label), err)
		exit.Error(err, "Failed to parse upsert body")
	}

	pc := sdk.NewPineconeClient(ctx)
	ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, options.namespace)
	if err != nil {
		msg.FailMsg("Failed to create index connection: %s", err)
		exit.Error(err, "Failed to create index connection")
	}

	if len(payload.Records) == 0 {
		msg.FailMsg("No records found in %s", style.Emphasis(src.Label))
		exit.ErrorMsg("No records provided for upsert")
	}

	records := make([]*pinecone.IntegratedRecord, 0, len(payload.Records))
	for i := range payload.Records {
		records = append(records, &payload.Records[i])
	}

	if options.batchSize <= 0 {
		options.batchSize = len(records)
	}

	batches := make([][]*pinecone.IntegratedRecord, 0, (len(records)+options.batchSize-1)/options.batchSize)
	for i := 0; i < len(records); i += options.batchSize {
		end := i + options.batchSize
		if end > len(records) {
			end = len(records)
		}
		batches = append(batches, records[i:end])
	}

	for i, batch := range batches {
		err := ic.UpsertRecords(ctx, batch)
		if err != nil {
			msg.FailMsg("Failed to upsert %d records in batch %d: %s", len(batch), i+1, err)
			exit.Errorf(err, "Failed to upsert %d records in batch %d", len(batch), i+1)
		} else if options.json {
			summary := map[string]any{
				"batch":     i + 1,
				"batches":   len(batches),
				"records":   len(batch),
				"namespace": options.namespace,
			}
			pcio.Println(text.IndentJSON(summary))
		} else {
			msg.SuccessMsg("Upserted %d records into namespace %s (batch %d of %d)", len(batch), options.namespace, i+1, len(batches))
		}
	}
}

func parseUpsertRecordsBody(b []byte) (*UpsertRecordsBody, error) {
	// First try JSON object: {"records":[...]}
	{
		var payload UpsertRecordsBody
		dec := json.NewDecoder(bytes.NewReader(b))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&payload); err == nil && len(payload.Records) > 0 {
			return &payload, nil
		}
	}

	// Fallback: JSONL/stream of pinecone.IntegratedRecord values
	var records []pinecone.IntegratedRecord
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	for {
		var rec pinecone.IntegratedRecord
		if err := dec.Decode(&rec); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	if len(records) == 0 {
		return nil, io.EOF
	}
	return &UpsertRecordsBody{Records: records}, nil
}
