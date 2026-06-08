package importcmd

import (
	"context"
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_runCancelImportCmd_RequiresID(t *testing.T) {
	svc := &mockImportService{}
	opts := cancelImportCmdOptions{}

	err := runCancelImportCmd(context.Background(), svc, opts)

	assert.Error(t, err)
	assert.Empty(t, svc.lastCancelImportId)
}

func Test_runCancelImportCmd_Succeeds(t *testing.T) {
	svc := &mockImportService{}
	opts := cancelImportCmdOptions{importId: "import-1"}

	err := runCancelImportCmd(context.Background(), svc, opts)

	if assert.NoError(t, err) {
		assert.Equal(t, "import-1", svc.lastCancelImportId)
	}
}

func Test_runCancelImportCmd_SucceedsJSON(t *testing.T) {
	svc := &mockImportService{}
	opts := cancelImportCmdOptions{importId: "import-1", json: true}

	out := testutils.CaptureStdout(t, func() {
		err := runCancelImportCmd(context.Background(), svc, opts)
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"import-1"`)
	assert.Contains(t, out, `"cancelled": true`)
}

func Test_runCancelImportCmd_PropagatesError(t *testing.T) {
	svc := &mockImportService{
		cancelImportErr: errors.New("cancel failed"),
	}
	opts := cancelImportCmdOptions{importId: "import-1"}

	err := runCancelImportCmd(context.Background(), svc, opts)

	assert.EqualError(t, err, "cancel failed")
}
