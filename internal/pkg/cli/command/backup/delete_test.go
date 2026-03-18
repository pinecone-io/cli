package backup

import (
	"context"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_runDeleteBackupCmd_RequiresBackupId(t *testing.T) {
	svc := &mockBackupService{}
	opts := deleteBackupCmdOptions{}

	err := runDeleteBackupCmd(context.Background(), svc, opts)

	assert.Error(t, err)
	assert.Empty(t, svc.lastDeleteBackupId)
}

func Test_runDeleteBackupCmd_Succeeds(t *testing.T) {
	svc := &mockBackupService{}
	opts := deleteBackupCmdOptions{backupId: "b1"}

	err := runDeleteBackupCmd(context.Background(), svc, opts)

	assert.NoError(t, err)
	assert.Equal(t, "b1", svc.lastDeleteBackupId)
}

func Test_runDeleteBackupCmd_SucceedsJSON(t *testing.T) {
	svc := &mockBackupService{}
	opts := deleteBackupCmdOptions{backupId: "b1", json: true}

	out := testutils.CaptureStdout(t, func() {
		err := runDeleteBackupCmd(context.Background(), svc, opts)
		assert.NoError(t, err)
	})

	assert.JSONEq(t, `{"deleted":true,"id":"b1"}`, out)
}
