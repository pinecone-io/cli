package organization

import (
	"context"
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

type mockDeleteOrganizationService struct {
	lastDeleteId string
	deleteErr    error
}

func (m *mockDeleteOrganizationService) Delete(ctx context.Context, id string) error {
	m.lastDeleteId = id
	return m.deleteErr
}

func Test_runDeleteOrganizationCmd_Succeeds(t *testing.T) {
	svc := &mockDeleteOrganizationService{}
	opts := deleteOrganizationCmdOptions{organizationID: "org-123"}

	err := runDeleteOrganizationCmd(context.Background(), svc, opts, "my-org", "org-123")

	assert.NoError(t, err)
	assert.Equal(t, "org-123", svc.lastDeleteId)
}

func Test_runDeleteOrganizationCmd_SucceedsJSON(t *testing.T) {
	svc := &mockDeleteOrganizationService{}
	opts := deleteOrganizationCmdOptions{organizationID: "org-123", json: true}

	out := testutils.CaptureStdout(t, func() {
		err := runDeleteOrganizationCmd(context.Background(), svc, opts, "my-org", "org-123")
		assert.NoError(t, err)
	})

	assert.JSONEq(t, `{"deleted":true,"name":"my-org","id":"org-123"}`, out)
}

func Test_runDeleteOrganizationCmd_PropagatesError(t *testing.T) {
	svc := &mockDeleteOrganizationService{deleteErr: errors.New("service error")}
	opts := deleteOrganizationCmdOptions{organizationID: "org-123"}

	err := runDeleteOrganizationCmd(context.Background(), svc, opts, "my-org", "org-123")

	assert.Error(t, err)
}
