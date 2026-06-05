package importcmd

import (
	"context"
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/stretchr/testify/assert"
)

func Test_runStartImportCmd_RequiresURI(t *testing.T) {
	svc := &mockImportService{}
	opts := startImportCmdOptions{}

	err := runStartImportCmd(context.Background(), svc, opts)

	assert.ErrorContains(t, err, "--uri is required")
}

func Test_runStartImportCmd_Succeeds(t *testing.T) {
	svc := &mockImportService{
		startImportResp: &pinecone.StartImportResponse{Id: "import-1"},
	}
	opts := startImportCmdOptions{
		uri: "s3://my-bucket/data/",
	}

	err := runStartImportCmd(context.Background(), svc, opts)

	if assert.NoError(t, err) {
		assert.Equal(t, "s3://my-bucket/data/", svc.lastStartImportUri)
	}
}

func Test_runStartImportCmd_PassesOptionalFields(t *testing.T) {
	svc := &mockImportService{
		startImportResp: &pinecone.StartImportResponse{Id: "import-1"},
	}
	opts := startImportCmdOptions{
		uri:           "s3://my-bucket/data/",
		integrationId: "intg-123",
		errorMode:     "abort",
	}

	err := runStartImportCmd(context.Background(), svc, opts)

	if assert.NoError(t, err) {
		if assert.NotNil(t, svc.lastStartImportIntegrationId) {
			assert.Equal(t, "intg-123", *svc.lastStartImportIntegrationId)
		}
		if assert.NotNil(t, svc.lastStartImportErrorMode) {
			assert.Equal(t, "abort", *svc.lastStartImportErrorMode)
		}
		assert.Equal(t, "s3://my-bucket/data/", svc.lastStartImportUri)

	}
}

func Test_runStartImportCmd_SucceedsJSON(t *testing.T) {
	svc := &mockImportService{
		startImportResp: &pinecone.StartImportResponse{Id: "import-1"},
	}
	opts := startImportCmdOptions{
		uri:  "s3://my-bucket/data/",
		json: true,
	}

	out := testutils.CaptureStdout(t, func() {
		err := runStartImportCmd(context.Background(), svc, opts)
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"import-1"`)
}

func Test_runStartImportCmd_PropagatesError(t *testing.T) {
	svc := &mockImportService{
		startImportErr: errors.New("start failed"),
	}
	opts := startImportCmdOptions{uri: "s3://my-bucket/data/"}

	err := runStartImportCmd(context.Background(), svc, opts)

	assert.EqualError(t, err, "start failed")
}
