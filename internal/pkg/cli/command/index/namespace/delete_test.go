package namespace

import (
	"context"
	"testing"

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
