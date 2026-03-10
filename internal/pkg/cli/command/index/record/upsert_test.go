package record

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// validation tests
// ---------------------------------------------------------------------------

func Test_runUpsertCmd_RequiresFile(t *testing.T) {
	svc := &mockRecordService{}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		batchSize: 96,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "--file or --body must be provided")
	assert.Empty(t, svc.upsertCalls)
}

func Test_runUpsertCmd_RejectsEmptyJSONArray(t *testing.T) {
	svc := &mockRecordService{}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		file:      `[]`,
		indexName: "my-index",
		namespace: "__default__",
		batchSize: 96,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no records provided")
	assert.Empty(t, svc.upsertCalls)
}

func Test_runUpsertCmd_RejectsEmptyJSONObject(t *testing.T) {
	svc := &mockRecordService{}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		file:      `{"records":[]}`,
		indexName: "my-index",
		namespace: "__default__",
		batchSize: 96,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no records provided")
	assert.Empty(t, svc.upsertCalls)
}

func Test_runUpsertCmd_RejectsInvalidJSON(t *testing.T) {
	svc := &mockRecordService{}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		file:      `not valid json at all`,
		indexName: "my-index",
		namespace: "__default__",
		batchSize: 96,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse upsert body")
	assert.Empty(t, svc.upsertCalls)
}

// ---------------------------------------------------------------------------
// success – input format variants
// ---------------------------------------------------------------------------

func Test_runUpsertCmd_JSONObjectFormat(t *testing.T) {
	svc := &mockRecordService{}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		file:      `{"records":[{"_id":"r1","chunk_text":"hello"},{"_id":"r2","chunk_text":"world"}]}`,
		indexName: "my-index",
		namespace: "my-ns",
		batchSize: 96,
	})

	require.NoError(t, err)
	require.Len(t, svc.upsertCalls, 1)
	assert.Len(t, svc.upsertCalls[0], 2)
}

func Test_runUpsertCmd_JSONArrayFormat(t *testing.T) {
	svc := &mockRecordService{}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		file:      `[{"_id":"r1","chunk_text":"hello"},{"_id":"r2","chunk_text":"world"}]`,
		indexName: "my-index",
		namespace: "__default__",
		batchSize: 96,
	})

	require.NoError(t, err)
	require.Len(t, svc.upsertCalls, 1)
	assert.Len(t, svc.upsertCalls[0], 2)
}

func Test_runUpsertCmd_JSONLFormat(t *testing.T) {
	jsonl := "{\"_id\":\"r1\",\"chunk_text\":\"hello\"}\n{\"_id\":\"r2\",\"chunk_text\":\"world\"}\n"
	svc := &mockRecordService{}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		file:      jsonl,
		indexName: "my-index",
		namespace: "__default__",
		batchSize: 96,
	})

	require.NoError(t, err)
	require.Len(t, svc.upsertCalls, 1)
	assert.Len(t, svc.upsertCalls[0], 2)
}

// ---------------------------------------------------------------------------
// batching
// ---------------------------------------------------------------------------

func Test_runUpsertCmd_BatchesByBatchSize(t *testing.T) {
	svc := &mockRecordService{}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		file:      `{"records":[{"_id":"r1"},{"_id":"r2"},{"_id":"r3"},{"_id":"r4"},{"_id":"r5"}]}`,
		indexName: "my-index",
		namespace: "__default__",
		batchSize: 2,
	})

	require.NoError(t, err)
	require.Len(t, svc.upsertCalls, 3, "5 records with batchSize=2 should produce 3 batches")
	assert.Len(t, svc.upsertCalls[0], 2)
	assert.Len(t, svc.upsertCalls[1], 2)
	assert.Len(t, svc.upsertCalls[2], 1)
}

func Test_runUpsertCmd_ZeroBatchSizeUpsertsAll(t *testing.T) {
	svc := &mockRecordService{}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		file:      `{"records":[{"_id":"r1"},{"_id":"r2"},{"_id":"r3"}]}`,
		indexName: "my-index",
		namespace: "__default__",
		batchSize: 0,
	})

	require.NoError(t, err)
	require.Len(t, svc.upsertCalls, 1, "batchSize=0 should upsert all records in a single call")
	assert.Len(t, svc.upsertCalls[0], 3)
}

// ---------------------------------------------------------------------------
// SDK error propagation
// ---------------------------------------------------------------------------

func Test_runUpsertCmd_PropagatesUpsertError(t *testing.T) {
	sdkErr := errors.New("upsert RPC failed")
	svc := &mockRecordService{upsertErr: sdkErr}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		file:      `{"records":[{"_id":"r1"}]}`,
		indexName: "my-index",
		namespace: "__default__",
		batchSize: 96,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, sdkErr)
	assert.Contains(t, err.Error(), "failed to upsert")
}

func Test_runUpsertCmd_StopsOnFirstBatchError(t *testing.T) {
	sdkErr := errors.New("rpc error")
	svc := &mockRecordService{upsertErr: sdkErr}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		file:      `{"records":[{"_id":"r1"},{"_id":"r2"},{"_id":"r3"}]}`,
		indexName: "my-index",
		namespace: "__default__",
		batchSize: 1,
	})

	assert.Error(t, err)
	assert.Len(t, svc.upsertCalls, 1, "should stop after the first failing batch")
}

// ---------------------------------------------------------------------------
// JSON output mode
// ---------------------------------------------------------------------------

func Test_runUpsertCmd_JSONOutput(t *testing.T) {
	svc := &mockRecordService{}
	err := runUpsertCmd(context.Background(), svc, upsertCmdOptions{
		file:      `{"records":[{"_id":"r1"},{"_id":"r2"}]}`,
		indexName: "my-index",
		namespace: "my-ns",
		batchSize: 96,
		json:      true,
	})

	require.NoError(t, err)
	require.Len(t, svc.upsertCalls, 1)
}
