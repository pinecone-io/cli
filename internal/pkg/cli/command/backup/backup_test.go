package backup

import (
	"context"
	"os"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/stretchr/testify/assert"
)

type mockBackupService struct {
	lastCreateBackupReq          *pinecone.CreateBackupParams
	lastDescribeBackupId         string
	lastListBackupsParams        *pinecone.ListBackupsParams
	lastDeleteBackupId           string
	lastCreateIndexFromBackupReq *pinecone.CreateIndexFromBackupParams

	createBackupResp          *pinecone.Backup
	describeBackupResp        *pinecone.Backup
	listBackupsResp           *pinecone.BackupList
	createIndexFromBackupResp *pinecone.CreateIndexFromBackupResponse

	createBackupErr          error
	describeBackupErr        error
	listBackupsErr           error
	deleteBackupErr          error
	createIndexFromBackupErr error
}

func (m *mockBackupService) CreateBackup(ctx context.Context, in *pinecone.CreateBackupParams) (*pinecone.Backup, error) {
	m.lastCreateBackupReq = in
	return m.createBackupResp, m.createBackupErr
}

func (m *mockBackupService) DescribeBackup(ctx context.Context, backupId string) (*pinecone.Backup, error) {
	m.lastDescribeBackupId = backupId
	return m.describeBackupResp, m.describeBackupErr
}

func (m *mockBackupService) ListBackups(ctx context.Context, in *pinecone.ListBackupsParams) (*pinecone.BackupList, error) {
	m.lastListBackupsParams = in
	return m.listBackupsResp, m.listBackupsErr
}

func (m *mockBackupService) DeleteBackup(ctx context.Context, backupId string) error {
	m.lastDeleteBackupId = backupId
	return m.deleteBackupErr
}

func (m *mockBackupService) CreateIndexFromBackup(ctx context.Context, in *pinecone.CreateIndexFromBackupParams) (*pinecone.CreateIndexFromBackupResponse, error) {
	m.lastCreateIndexFromBackupReq = in
	return m.createIndexFromBackupResp, m.createIndexFromBackupErr
}

func TestMain(m *testing.M) {
	reset := testutils.SilenceOutput()
	code := m.Run()
	reset()
	os.Exit(code)
}

func Test_runCreateBackupCmd_RequiresIndexName(t *testing.T) {
	svc := &mockBackupService{}
	opts := createBackupCmdOptions{}

	err := runCreateBackupCmd(context.Background(), svc, opts)

	assert.Error(t, err)
	assert.Nil(t, svc.lastCreateBackupReq)
}

func Test_runCreateBackupCmd_PopulatesRequest(t *testing.T) {
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
	opts := describeBackupCmdOptions{
		backupId: "b1",
	}

	err := runDescribeBackupCmd(context.Background(), svc, opts)

	assert.NoError(t, err)
	assert.Equal(t, "b1", svc.lastDescribeBackupId)
}

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
