package namespace

import (
	"context"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_runDeleteNamespaceCmd_RequiresName(t *testing.T) {
	svc := &mockNamespaceService{}
	options := deleteNamespaceCmdOptions{
		name: "",
	}

	err := runDeleteNamespaceCmd(context.Background(), svc, options)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--name is required")
	assert.Empty(t, svc.lastDeleteArg)
}

func Test_runDeleteNamespaceCmd_Succeeds(t *testing.T) {
	svc := &mockNamespaceService{}
	options := deleteNamespaceCmdOptions{
		name: "tenant-a",
	}

	err := runDeleteNamespaceCmd(context.Background(), svc, options)

	assert.NoError(t, err)
	assert.Equal(t, options.name, svc.lastDeleteArg)
}

func Test_runDeleteNamespaceCmd_SucceedsJSON(t *testing.T) {
	svc := &mockNamespaceService{}
	options := deleteNamespaceCmdOptions{
		name:      "tenant-a",
		indexName: "my-index",
		json:      true,
	}

	out := testutils.CaptureStdout(t, func() {
		err := runDeleteNamespaceCmd(context.Background(), svc, options)
		assert.NoError(t, err)
	})

	assert.JSONEq(t, `{"deleted":true,"namespace":"tenant-a","index":"my-index"}`, out)
}
