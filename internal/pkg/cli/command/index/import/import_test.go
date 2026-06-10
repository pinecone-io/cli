package importcmd

import (
	"context"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

type mockImportService struct {
	lastStartImportUri           string
	lastStartImportIntegrationId *string
	lastStartImportErrorMode     *string
	lastDescribeImportId         string
	lastListImportsLimit         *int32
	lastListImportsPaginationToken *string
	lastCancelImportId           string

	startImportResp  *pinecone.StartImportResponse
	describeImportResp *pinecone.Import
	listImportsResp  *pinecone.ListImportsResponse

	startImportErr   error
	describeImportErr error
	listImportsErr   error
	cancelImportErr  error
}

func (m *mockImportService) StartImport(ctx context.Context, uri string, integrationId, errorMode *string) (*pinecone.StartImportResponse, error) {
	m.lastStartImportUri = uri
	m.lastStartImportIntegrationId = integrationId
	m.lastStartImportErrorMode = errorMode
	return m.startImportResp, m.startImportErr
}

func (m *mockImportService) DescribeImport(ctx context.Context, id string) (*pinecone.Import, error) {
	m.lastDescribeImportId = id
	return m.describeImportResp, m.describeImportErr
}

func (m *mockImportService) ListImports(ctx context.Context, limit *int32, paginationToken *string) (*pinecone.ListImportsResponse, error) {
	m.lastListImportsLimit = limit
	m.lastListImportsPaginationToken = paginationToken
	return m.listImportsResp, m.listImportsErr
}

func (m *mockImportService) CancelImport(ctx context.Context, id string) error {
	m.lastCancelImportId = id
	return m.cancelImportErr
}
