package backup

import (
	"context"

	"github.com/pinecone-io/go-pinecone/v5/pinecone"
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
