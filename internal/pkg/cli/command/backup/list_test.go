package backup

import (
	"context"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/stretchr/testify/assert"
)

func Test_runListBackupsCmd_PopulatesParams(t *testing.T) {
	svc := &mockBackupService{
		listBackupsResp: &pinecone.BackupList{},
	}
	opts := listBackupsCmdOptions{
		indexName:       "idx",
		limit:           5,
		paginationToken: "next",
	}

	err := runListBackupsCmd(context.Background(), svc, opts)

	if assert.NoError(t, err) {
		if assert.NotNil(t, svc.lastListBackupsParams) {
			if assert.NotNil(t, svc.lastListBackupsParams.IndexName) {
				assert.Equal(t, "idx", *svc.lastListBackupsParams.IndexName)
			}
			if assert.NotNil(t, svc.lastListBackupsParams.Limit) {
				assert.Equal(t, 5, *svc.lastListBackupsParams.Limit)
			}
			if assert.NotNil(t, svc.lastListBackupsParams.PaginationToken) {
				assert.Equal(t, "next", *svc.lastListBackupsParams.PaginationToken)
			}
		}
	}
}

func Test_runListBackupsCmd_SucceedsJSON(t *testing.T) {
	svc := &mockBackupService{
		listBackupsResp: &pinecone.BackupList{
			Data: []*pinecone.Backup{{BackupId: "b1"}},
		},
	}
	opts := listBackupsCmdOptions{json: true}

	out := testutils.CaptureStdout(t, func() {
		err := runListBackupsCmd(context.Background(), svc, opts)
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"b1"`)
}
