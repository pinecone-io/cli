package backup

import (
	"context"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/stretchr/testify/assert"
)

func Test_runCreateBackupCmd_RequiresIndexName(t *testing.T) {
	svc := &mockBackupService{}
	opts := createBackupCmdOptions{}

	err := runCreateBackupCmd(context.Background(), svc, opts)

	assert.Error(t, err)
	assert.Nil(t, svc.lastCreateBackupReq)
}

func Test_runCreateBackupCmd_Succeeds(t *testing.T) {
	svc := &mockBackupService{
		createBackupResp: &pinecone.Backup{BackupId: "b1"},
	}
	opts := createBackupCmdOptions{
		indexName:   "idx",
		description: "desc",
		name:        "name",
	}

	err := runCreateBackupCmd(context.Background(), svc, opts)

	if assert.NoError(t, err) {
		if assert.NotNil(t, svc.lastCreateBackupReq) {
			assert.Equal(t, "idx", svc.lastCreateBackupReq.IndexName)
			if assert.NotNil(t, svc.lastCreateBackupReq.Description) {
				assert.Equal(t, "desc", *svc.lastCreateBackupReq.Description)
			}
			if assert.NotNil(t, svc.lastCreateBackupReq.Name) {
				assert.Equal(t, "name", *svc.lastCreateBackupReq.Name)
			}
		}
	}
}

func Test_runCreateBackupCmd_SucceedsJSON(t *testing.T) {
	svc := &mockBackupService{
		createBackupResp: &pinecone.Backup{BackupId: "b1"},
	}
	opts := createBackupCmdOptions{
		indexName: "idx",
		json:      true,
	}

	out := testutils.CaptureStdout(t, func() {
		err := runCreateBackupCmd(context.Background(), svc, opts)
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"b1"`)
}
