package project

import (
	"context"
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

type mockDeleteProjectService struct {
	lastDeleteId string
	deleteErr    error
}

func (m *mockDeleteProjectService) Delete(ctx context.Context, id string) error {
	m.lastDeleteId = id
	return m.deleteErr
}

func Test_runDeleteProjectCmd_Succeeds(t *testing.T) {
	svc := &mockDeleteProjectService{}
	opts := deleteProjectCmdOptions{projectId: "proj-123"}

	err := runDeleteProjectCmd(context.Background(), svc, opts, "my-project", "proj-123")

	assert.NoError(t, err)
	assert.Equal(t, "proj-123", svc.lastDeleteId)
}

func Test_runDeleteProjectCmd_SucceedsJSON(t *testing.T) {
	svc := &mockDeleteProjectService{}
	opts := deleteProjectCmdOptions{projectId: "proj-123", json: true}

	out := testutils.CaptureStdout(t, func() {
		err := runDeleteProjectCmd(context.Background(), svc, opts, "my-project", "proj-123")
		assert.NoError(t, err)
	})

	assert.JSONEq(t, `{"deleted":true,"name":"my-project","id":"proj-123"}`, out)
}

func Test_runDeleteProjectCmd_PropagatesError(t *testing.T) {
	svc := &mockDeleteProjectService{deleteErr: errors.New("service error")}
	opts := deleteProjectCmdOptions{projectId: "proj-123"}

	err := runDeleteProjectCmd(context.Background(), svc, opts, "my-project", "proj-123")

	assert.Error(t, err)
}
