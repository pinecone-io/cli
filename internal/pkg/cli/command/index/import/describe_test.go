package importcmd

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/stretchr/testify/assert"
)

func Test_runDescribeImportCmd_RequiresID(t *testing.T) {
	svc := &mockImportService{}
	opts := describeImportCmdOptions{}

	err := runDescribeImportCmd(context.Background(), svc, opts)

	assert.Error(t, err)
	assert.Empty(t, svc.lastDescribeImportId)
}

func Test_runDescribeImportCmd_Succeeds(t *testing.T) {
	now := time.Now()
	svc := &mockImportService{
		describeImportResp: &pinecone.Import{
			Id:        "import-1",
			Status:    pinecone.InProgress,
			Uri:       "s3://my-bucket/data/",
			CreatedAt: &now,
		},
	}
	opts := describeImportCmdOptions{importId: "import-1"}

	err := runDescribeImportCmd(context.Background(), svc, opts)

	if assert.NoError(t, err) {
		assert.Equal(t, "import-1", svc.lastDescribeImportId)
	}
}

func Test_runDescribeImportCmd_SucceedsJSON(t *testing.T) {
	svc := &mockImportService{
		describeImportResp: &pinecone.Import{
			Id:     "import-1",
			Status: pinecone.Completed,
		},
	}
	opts := describeImportCmdOptions{importId: "import-1", json: true}

	out := testutils.CaptureStdout(t, func() {
		err := runDescribeImportCmd(context.Background(), svc, opts)
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"import-1"`)
}

func Test_runDescribeImportCmd_PropagatesError(t *testing.T) {
	svc := &mockImportService{
		describeImportErr: errors.New("not found"),
	}
	opts := describeImportCmdOptions{importId: "import-1"}

	err := runDescribeImportCmd(context.Background(), svc, opts)

	assert.EqualError(t, err, "not found")
}
