package namespace

import (
	"context"
	"os"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

var (
	mockNamespaceDescription = pinecone.NamespaceDescription{
		Name:        "tenant-a",
		RecordCount: 42,
	}
	mockListNamespacesResponse = pinecone.ListNamespacesResponse{
		Namespaces: []*pinecone.NamespaceDescription{
			{
				Name:        "tenant-a",
				RecordCount: 1,
			},
		},
		TotalCount: 1,
	}
)

type mockNamespaceService struct {
	lastCreateReq   *pinecone.CreateNamespaceParams
	lastDescribeArg string
	lastListParams  *pinecone.ListNamespacesParams
	lastDeleteArg   string

	createResp   *pinecone.NamespaceDescription
	describeResp *pinecone.NamespaceDescription
	listResp     *pinecone.ListNamespacesResponse

	createErr   error
	describeErr error
	listErr     error
	deleteErr   error
}

func (m *mockNamespaceService) CreateNamespace(ctx context.Context, req *pinecone.CreateNamespaceParams) (*pinecone.NamespaceDescription, error) {
	m.lastCreateReq = req
	return m.createResp, m.createErr
}

func (m *mockNamespaceService) DescribeNamespace(ctx context.Context, name string) (*pinecone.NamespaceDescription, error) {
	m.lastDescribeArg = name
	return m.describeResp, m.describeErr
}

func (m *mockNamespaceService) ListNamespaces(ctx context.Context, params *pinecone.ListNamespacesParams) (*pinecone.ListNamespacesResponse, error) {
	m.lastListParams = params
	return m.listResp, m.listErr
}

func (m *mockNamespaceService) DeleteNamespace(ctx context.Context, name string) error {
	m.lastDeleteArg = name
	return m.deleteErr
}

func TestMain(m *testing.M) {
	reset := testutils.SilenceOutput()
	code := m.Run()
	reset()
	os.Exit(code)
}
