package record

import (
	"context"
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/utils/flags"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// emptySearchResp is a non-nil response so that output helpers don't panic.
var emptySearchResp = &pinecone.SearchRecordsResponse{}

// ---------------------------------------------------------------------------
// validation
// ---------------------------------------------------------------------------

func Test_runSearchCmd_RequiresQueryMode(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      10,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "provide a query")
	assert.Nil(t, svc.lastSearchReq)
}

func Test_runSearchCmd_RejectsZeroTopK(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	id := "rec-1"
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      0,
		id:        id,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "top-k must be greater than 0")
	assert.Nil(t, svc.lastSearchReq)
}

func Test_runSearchCmd_RejectsSparseLengthMismatch(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName:     "my-index",
		namespace:     "__default__",
		topK:          10,
		sparseIndices: flags.Int32List{1, 2},
		sparseValues:  flags.Float32List{0.5}, // wrong length
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "same length")
	assert.Nil(t, svc.lastSearchReq)
}

// ---------------------------------------------------------------------------
// query modes
// ---------------------------------------------------------------------------

func Test_runSearchCmd_IDQuery(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      10,
		id:        "rec-123",
	})

	require.NoError(t, err)
	require.NotNil(t, svc.lastSearchReq)
	require.NotNil(t, svc.lastSearchReq.Query.Id)
	assert.Equal(t, "rec-123", *svc.lastSearchReq.Query.Id)
	assert.Equal(t, int32(10), svc.lastSearchReq.Query.TopK)
}

func Test_runSearchCmd_InputsQuery(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      5,
		inputs:    flags.JSONObject{"text": "disease prevention"},
	})

	require.NoError(t, err)
	require.NotNil(t, svc.lastSearchReq)
	require.NotNil(t, svc.lastSearchReq.Query.Inputs)
	assert.Equal(t, "disease prevention", (*svc.lastSearchReq.Query.Inputs)["text"])
	assert.Equal(t, int32(5), svc.lastSearchReq.Query.TopK)
}

func Test_runSearchCmd_DenseVectorQuery(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      10,
		vector:    flags.Float32List{0.1, 0.2, 0.3},
	})

	require.NoError(t, err)
	require.NotNil(t, svc.lastSearchReq)
	require.NotNil(t, svc.lastSearchReq.Query.Vector)
	require.NotNil(t, svc.lastSearchReq.Query.Vector.Values)
	assert.Equal(t, []float32{0.1, 0.2, 0.3}, *svc.lastSearchReq.Query.Vector.Values)
	assert.Nil(t, svc.lastSearchReq.Query.Vector.SparseIndices)
}

func Test_runSearchCmd_SparseVectorQuery(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName:     "my-index",
		namespace:     "__default__",
		topK:          10,
		sparseIndices: flags.Int32List{1, 5, 9},
		sparseValues:  flags.Float32List{0.4, 0.8, 0.2},
	})

	require.NoError(t, err)
	require.NotNil(t, svc.lastSearchReq)
	require.NotNil(t, svc.lastSearchReq.Query.Vector)
	require.NotNil(t, svc.lastSearchReq.Query.Vector.SparseIndices)
	assert.Equal(t, []int32{1, 5, 9}, *svc.lastSearchReq.Query.Vector.SparseIndices)
	assert.Equal(t, []float32{0.4, 0.8, 0.2}, *svc.lastSearchReq.Query.Vector.SparseValues)
	assert.Nil(t, svc.lastSearchReq.Query.Vector.Values)
}

func Test_runSearchCmd_HybridQuery(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName:     "my-index",
		namespace:     "__default__",
		topK:          10,
		vector:        flags.Float32List{0.1, 0.2, 0.3},
		sparseIndices: flags.Int32List{1, 5},
		sparseValues:  flags.Float32List{0.4, 0.8},
	})

	require.NoError(t, err)
	require.NotNil(t, svc.lastSearchReq.Query.Vector)
	assert.NotNil(t, svc.lastSearchReq.Query.Vector.Values)
	assert.NotNil(t, svc.lastSearchReq.Query.Vector.SparseIndices)
}

// ---------------------------------------------------------------------------
// optional modifiers
// ---------------------------------------------------------------------------

func Test_runSearchCmd_AppliesFilter(t *testing.T) {
	wantFilter := map[string]interface{}{"category": "health"}
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      10,
		id:        "rec-1",
		filter:    flags.JSONObject(wantFilter),
	})

	require.NoError(t, err)
	require.NotNil(t, svc.lastSearchReq.Query.Filter)
	assert.Equal(t, wantFilter, *svc.lastSearchReq.Query.Filter)
}

func Test_runSearchCmd_AppliesFields(t *testing.T) {
	wantFields := []string{"_id", "chunk_text"}
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      10,
		id:        "rec-1",
		fields:    flags.StringList(wantFields),
	})

	require.NoError(t, err)
	require.NotNil(t, svc.lastSearchReq.Fields)
	assert.Equal(t, wantFields, *svc.lastSearchReq.Fields)
}

func Test_runSearchCmd_AppliesRerank(t *testing.T) {
	wantRerank := pinecone.SearchRecordsRerank{
		Model:      "bge-reranker-v2-m3",
		RankFields: []string{"chunk_text"},
	}
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      10,
		id:        "rec-1",
		rerank:    flags.JSONObject{"model": wantRerank.Model, "rank_fields": []interface{}{wantRerank.RankFields[0]}},
	})

	require.NoError(t, err)
	require.NotNil(t, svc.lastSearchReq.Rerank)
	assert.Equal(t, wantRerank, *svc.lastSearchReq.Rerank)
}

func Test_runSearchCmd_AppliesMatchTerms(t *testing.T) {
	wantTerms := []string{"vaccine", "prevention"}
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName:  "my-index",
		namespace:  "__default__",
		topK:       10,
		inputs:     flags.JSONObject{"text": "disease"},
		matchTerms: flags.JSONObject{"terms": []interface{}{"vaccine", "prevention"}},
	})

	require.NoError(t, err)
	require.NotNil(t, svc.lastSearchReq.Query.MatchTerms)
	require.NotNil(t, svc.lastSearchReq.Query.MatchTerms.Terms)
	assert.Equal(t, wantTerms, *svc.lastSearchReq.Query.MatchTerms.Terms)
}

// ---------------------------------------------------------------------------
// --body overlay and flag precedence
// ---------------------------------------------------------------------------

func Test_runSearchCmd_BodyProvidesQuery(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      defaultSearchTopK,
		// No query flags set; body supplies the id.
		body: `{"query":{"top_k":3,"id":"body-rec"}}`,
	})

	require.NoError(t, err)
	require.NotNil(t, svc.lastSearchReq.Query.Id)
	assert.Equal(t, "body-rec", *svc.lastSearchReq.Query.Id)
	assert.Equal(t, int32(3), svc.lastSearchReq.Query.TopK)
}

func Test_runSearchCmd_FlagIdWinsOverBody(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      defaultSearchTopK,
		id:        "flag-rec",
		body:      `{"query":{"id":"body-rec"}}`,
	})

	require.NoError(t, err)
	require.NotNil(t, svc.lastSearchReq.Query.Id)
	assert.Equal(t, "flag-rec", *svc.lastSearchReq.Query.Id)
}

func Test_runSearchCmd_ExplicitTopKWinsOverBody(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName:    "my-index",
		namespace:    "__default__",
		topK:         20,
		topKExplicit: true,
		id:           "rec-1",
		body:         `{"query":{"top_k":99,"id":"rec-1"}}`,
	})

	require.NoError(t, err)
	assert.Equal(t, int32(20), svc.lastSearchReq.Query.TopK)
}

func Test_runSearchCmd_BodyTopKAppliedWhenNotExplicit(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName:    "my-index",
		namespace:    "__default__",
		topK:         defaultSearchTopK, // default, not explicitly set
		topKExplicit: false,
		body:         `{"query":{"top_k":7,"id":"rec-1"}}`,
	})

	require.NoError(t, err)
	assert.Equal(t, int32(7), svc.lastSearchReq.Query.TopK)
}

// ---------------------------------------------------------------------------
// SDK error propagation and output
// ---------------------------------------------------------------------------

func Test_runSearchCmd_PropagatesSDKError(t *testing.T) {
	sdkErr := errors.New("search RPC failed")
	svc := &mockRecordService{searchErr: sdkErr}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      10,
		id:        "rec-1",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, sdkErr)
}

func Test_runSearchCmd_JSONOutput(t *testing.T) {
	svc := &mockRecordService{searchResp: emptySearchResp}
	err := runSearchCmd(context.Background(), svc, searchCmdOptions{
		indexName: "my-index",
		namespace: "__default__",
		topK:      10,
		id:        "rec-1",
		json:      true,
	})

	require.NoError(t, err)
	assert.NotNil(t, svc.lastSearchReq)
}
