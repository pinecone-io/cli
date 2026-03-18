package index

import (
	"context"
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

type mockDeleteIndexService struct {
	lastDeleteArg string
	deleteErr     error
}

func (m *mockDeleteIndexService) DeleteIndex(ctx context.Context, name string) error {
	m.lastDeleteArg = name
	return m.deleteErr
}

func Test_runDeleteIndexCmd_Succeeds(t *testing.T) {
	svc := &mockDeleteIndexService{}
	opts := deleteCmdOptions{name: "my-index"}

	err := runDeleteIndexCmd(context.Background(), svc, opts)

	assert.NoError(t, err)
	assert.Equal(t, "my-index", svc.lastDeleteArg)
}

func Test_runDeleteIndexCmd_SucceedsJSON(t *testing.T) {
	svc := &mockDeleteIndexService{}
	opts := deleteCmdOptions{name: "my-index", json: true}

	out := testutils.CaptureStdout(t, func() {
		err := runDeleteIndexCmd(context.Background(), svc, opts)
		assert.NoError(t, err)
	})

	assert.JSONEq(t, `{"deleted":true,"name":"my-index"}`, out)
}

func Test_runDeleteIndexCmd_PropagatesError(t *testing.T) {
	svc := &mockDeleteIndexService{deleteErr: errors.New("not found")}
	opts := deleteCmdOptions{name: "missing"}

	err := runDeleteIndexCmd(context.Background(), svc, opts)

	assert.Error(t, err)
	assert.Equal(t, "missing", svc.lastDeleteArg)
}
