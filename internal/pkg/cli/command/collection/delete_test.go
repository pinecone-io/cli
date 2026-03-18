package collection

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

type mockDeleteCollectionService struct {
	lastDeleteArg string
	deleteErr     error
}

func (m *mockDeleteCollectionService) DeleteCollection(ctx context.Context, name string) error {
	m.lastDeleteArg = name
	return m.deleteErr
}

func TestMain(m *testing.M) {
	reset := testutils.SilenceOutput()
	code := m.Run()
	reset()
	os.Exit(code)
}

func Test_runDeleteCollectionCmd_Succeeds(t *testing.T) {
	svc := &mockDeleteCollectionService{}
	opts := deleteCollectionCmdOptions{name: "my-collection"}

	err := runDeleteCollectionCmd(context.Background(), svc, opts)

	assert.NoError(t, err)
	assert.Equal(t, "my-collection", svc.lastDeleteArg)
}

func Test_runDeleteCollectionCmd_SucceedsJSON(t *testing.T) {
	svc := &mockDeleteCollectionService{}
	opts := deleteCollectionCmdOptions{name: "my-collection", json: true}

	out := testutils.CaptureStdout(t, func() {
		err := runDeleteCollectionCmd(context.Background(), svc, opts)
		assert.NoError(t, err)
	})

	assert.JSONEq(t, `{"deleted":true,"name":"my-collection"}`, out)
}

func Test_runDeleteCollectionCmd_PropagatesError(t *testing.T) {
	svc := &mockDeleteCollectionService{deleteErr: errors.New("service error")}
	opts := deleteCollectionCmdOptions{name: "my-collection"}

	err := runDeleteCollectionCmd(context.Background(), svc, opts)

	assert.Error(t, err)
}
