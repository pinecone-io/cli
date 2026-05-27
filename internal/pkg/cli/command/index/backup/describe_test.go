package backup

import (
	"context"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/stretchr/testify/assert"
)

func Test_runDescribeBackupCmd_RequiresBackupId(t *testing.T) {
	svc := &mockBackupService{}
	opts := describeBackupCmdOptions{}

	err := runDescribeBackupCmd(context.Background(), svc, opts)

	assert.Error(t, err)
	assert.Empty(t, svc.lastDescribeBackupId)
}

func Test_runDescribeBackupCmd_Succeeds(t *testing.T) {
	svc := &mockBackupService{
		describeBackupResp: &pinecone.Backup{BackupId: "b1"},
	}
	opts := describeBackupCmdOptions{backupId: "b1"}

	err := runDescribeBackupCmd(context.Background(), svc, opts)

	assert.NoError(t, err)
	assert.Equal(t, "b1", svc.lastDescribeBackupId)
}

func Test_runDescribeBackupCmd_SucceedsJSON(t *testing.T) {
	svc := &mockBackupService{
		describeBackupResp: &pinecone.Backup{BackupId: "b1"},
	}
	opts := describeBackupCmdOptions{backupId: "b1", json: true}

	out := testutils.CaptureStdout(t, func() {
		err := runDescribeBackupCmd(context.Background(), svc, opts)
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"b1"`)
}
