package restore

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/stretchr/testify/assert"
)

type mockRestoreJobService struct {
	lastDescribeId               string
	lastListParams               *pinecone.ListRestoreJobsParams
	lastCreateIndexFromBackupReq *pinecone.CreateIndexFromBackupParams

	describeResp              *pinecone.RestoreJob
	listResp                  *pinecone.RestoreJobList
	createIndexFromBackupResp *pinecone.CreateIndexFromBackupResponse

	describeErr              error
	listErr                  error
	createIndexFromBackupErr error
}

func (m *mockRestoreJobService) DescribeRestoreJob(ctx context.Context, restoreJobId string) (*pinecone.RestoreJob, error) {
	m.lastDescribeId = restoreJobId
	return m.describeResp, m.describeErr
}

func (m *mockRestoreJobService) ListRestoreJobs(ctx context.Context, in *pinecone.ListRestoreJobsParams) (*pinecone.RestoreJobList, error) {
	m.lastListParams = in
	return m.listResp, m.listErr
}

func (m *mockRestoreJobService) CreateIndexFromBackup(ctx context.Context, in *pinecone.CreateIndexFromBackupParams) (*pinecone.CreateIndexFromBackupResponse, error) {
	m.lastCreateIndexFromBackupReq = in
	return m.createIndexFromBackupResp, m.createIndexFromBackupErr
}

func TestMain(m *testing.M) {
	reset := testutils.SilenceOutput()
	code := m.Run()
	reset()
	os.Exit(code)
}

func Test_runDescribeRestoreJobCmd_RequiresId(t *testing.T) {
	svc := &mockRestoreJobService{}
	opts := describeRestoreJobCmdOptions{}

	err := runDescribeRestoreJobCmd(context.Background(), svc, opts)

	assert.Error(t, err)
	assert.Empty(t, svc.lastDescribeId)
}

func Test_runDescribeRestoreJobCmd_Succeeds(t *testing.T) {
	now := time.Now()
	svc := &mockRestoreJobService{
		describeResp: &pinecone.RestoreJob{
			RestoreJobId:    "rj-1",
			BackupId:        "b1",
			TargetIndexName: "idx",
			Status:          "completed",
			CreatedAt:       now,
		},
	}
	opts := describeRestoreJobCmdOptions{restoreJobId: "rj-1"}

	err := runDescribeRestoreJobCmd(context.Background(), svc, opts)

	assert.NoError(t, err)
	assert.Equal(t, "rj-1", svc.lastDescribeId)
}

func Test_runListRestoreJobsCmd_PopulatesParams(t *testing.T) {
	svc := &mockRestoreJobService{
		listResp: &pinecone.RestoreJobList{},
	}
	opts := listRestoreJobsCmdOptions{
		limit:           10,
		paginationToken: "next",
	}

	err := runListRestoreJobsCmd(context.Background(), svc, opts)

	if assert.NoError(t, err) {
		if assert.NotNil(t, svc.lastListParams) {
			if assert.NotNil(t, svc.lastListParams.Limit) {
				assert.Equal(t, 10, *svc.lastListParams.Limit)
			}
			if assert.NotNil(t, svc.lastListParams.PaginationToken) {
				assert.Equal(t, "next", *svc.lastListParams.PaginationToken)
			}
		}
	}
}

func Test_runRestoreJobCmd_ValidatesRequired(t *testing.T) {
	svc := &mockRestoreJobService{}

	err := runRestoreJobCmd(context.Background(), svc, restoreJobCmdOptions{})
	assert.Error(t, err)

	err = runRestoreJobCmd(context.Background(), svc, restoreJobCmdOptions{backupId: "b1"})
	assert.Error(t, err)
}

func Test_runCreateIndexFromBackupCmd_Succeeds(t *testing.T) {
	svc := &mockRestoreJobService{
		createIndexFromBackupResp: &pinecone.CreateIndexFromBackupResponse{
			IndexId:      "idx-id",
			RestoreJobId: "rj-1",
		},
	}
	opts := restoreJobCmdOptions{
		backupId:           "b1",
		name:               "new-index",
		deletionProtection: "enabled",
		tags:               map[string]string{"env": "prod"},
	}

	err := runRestoreJobCmd(context.Background(), svc, opts)

	if assert.NoError(t, err) {
		if assert.NotNil(t, svc.lastCreateIndexFromBackupReq) {
			assert.Equal(t, "b1", svc.lastCreateIndexFromBackupReq.BackupId)
			assert.Equal(t, "new-index", svc.lastCreateIndexFromBackupReq.Name)
			if assert.NotNil(t, svc.lastCreateIndexFromBackupReq.DeletionProtection) {
				assert.Equal(t, pinecone.DeletionProtectionEnabled, *svc.lastCreateIndexFromBackupReq.DeletionProtection)
			}
			if assert.NotNil(t, svc.lastCreateIndexFromBackupReq.Tags) {
				assert.Equal(t, "prod", (*svc.lastCreateIndexFromBackupReq.Tags)["env"])
			}
		}
	}
}
