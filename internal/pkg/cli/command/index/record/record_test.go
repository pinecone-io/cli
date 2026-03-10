package record

import (
	"context"
	"os"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

// mockRecordService is the shared test double for both upsert and search tests.
type mockRecordService struct {
	// upsert
	upsertErr   error
	upsertCalls [][]*pinecone.IntegratedRecord

	// search
	searchResp    *pinecone.SearchRecordsResponse
	searchErr     error
	lastSearchReq *pinecone.SearchRecordsRequest
}

func (m *mockRecordService) UpsertRecords(_ context.Context, records []*pinecone.IntegratedRecord) error {
	m.upsertCalls = append(m.upsertCalls, records)
	return m.upsertErr
}

func (m *mockRecordService) SearchRecords(_ context.Context, req *pinecone.SearchRecordsRequest) (*pinecone.SearchRecordsResponse, error) {
	m.lastSearchReq = req
	return m.searchResp, m.searchErr
}

func TestMain(m *testing.M) {
	reset := testutils.SilenceOutput()
	code := m.Run()
	reset()
	os.Exit(code)
}
