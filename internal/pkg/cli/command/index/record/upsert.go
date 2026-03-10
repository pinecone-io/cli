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

// Ensure *pinecone.IndexConnection satisfies RecordService at compile time.
var _ RecordService = (*pinecone.IndexConnection)(nil)

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
		Short: "Upsert text records into an integrated index from a JSON/JSONL payload",
		Long: help.Long(`
			Upsert records into an integrated index namespace from a JSON or JSONL payload.

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
			ctx := cmd.Context()
			pc := sdk.NewPineconeClient(ctx)
			ic, err := sdk.NewIndexConnection(ctx, pc, options.indexName, options.namespace)
			if err != nil {
				msg.FailMsg("Failed to create index connection: %s", err)
				exit.Error(err, "Failed to create index connection")
			}
			if err := runUpsertCmd(ctx, ic, options); err != nil {
				msg.FailMsg("%s", err)
				exit.Error(err, "upsert failed")
			}
		},
	}

	cmd.Flags().StringVarP(&options.indexName, "index-name", "n", "", "name of index to upsert into")
	cmd.Flags().StringVar(&options.namespace, "namespace", "__default__", "namespace to upsert into")
	cmd.Flags().StringVar(&options.file, "file", "", "request body JSON or JSONL (inline, ./path.json[l], or '-' for stdin; only one argument may use stdin)")
	cmd.Flags().StringVar(&options.file, "body", "", "alias for --file")
	cmd.Flags().IntVarP(&options.batchSize, "batch-size", "b", 96, "records per batch (max 96)")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "output as JSON")

	_ = cmd.MarkFlagRequired("index-name")

	return cmd
}

func runUpsertCmd(ctx context.Context, ic RecordService, options upsertCmdOptions) error {
	if options.file == "" {
		return pcio.Errorf("either --file or --body must be provided")
	}

	b, src, err := argio.ReadAll(options.file)
	if err != nil {
		return pcio.Errorf("failed to read upsert body (%s): %w", style.Emphasis(src.Label), err)
	}

	payload, err := parseUpsertRecordsBody(b)
	if err != nil {
		return pcio.Errorf("failed to parse upsert body (%s): %w", style.Emphasis(src.Label), err)
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
		if err := ic.UpsertRecords(ctx, batch); err != nil {
			return pcio.Errorf("failed to upsert %d records in batch %d: %w", len(batch), i+1, err)
		}
		if options.json {
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

	return nil
}

func parseUpsertRecordsBody(b []byte) (*UpsertRecordsBody, error) {
	// First try JSON object: {"records":[...]}
	{
		var payload UpsertRecordsBody
		dec := json.NewDecoder(bytes.NewReader(b))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&payload); err == nil {
			if len(payload.Records) == 0 {
				return nil, pcio.Errorf("no records provided")
			}
			return &payload, nil
		}
	}

	// Second try: raw JSON array of IntegratedRecord objects
	{
		var records []pinecone.IntegratedRecord
		dec := json.NewDecoder(bytes.NewReader(b))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&records); err == nil {
			if len(records) == 0 {
				return nil, pcio.Errorf("no records provided")
			}
			return &UpsertRecordsBody{Records: records}, nil
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
		return nil, pcio.Errorf("no records provided")
	}
	return &UpsertRecordsBody{Records: records}, nil
}
