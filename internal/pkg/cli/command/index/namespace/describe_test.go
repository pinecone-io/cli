package namespace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_runDescribeNamespaceCmd_RequiresName(t *testing.T) {
	svc := &mockNamespaceService{}
	options := describeNamespaceCmdOptions{
		name: "",
	}

	err := runDescribeNamespaceCmd(context.Background(), svc, options)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--name is required")
	assert.Empty(t, svc.lastDescribeArg)
}

func Test_runDescribeNamespaceCmd_Succeeds(t *testing.T) {
	svc := &mockNamespaceService{
		describeResp: &mockNamespaceDescription,
	}
	options := describeNamespaceCmdOptions{
		name: "tenant-a",
	}

	err := runDescribeNamespaceCmd(context.Background(), svc, options)

	assert.NoError(t, err)
	assert.Equal(t, options.name, svc.lastDescribeArg)
}

func Test_runDescribeNamespaceCmd_SucceedsJSON(t *testing.T) {
	svc := &mockNamespaceService{
		describeResp: &mockNamespaceDescription,
	}
	options := describeNamespaceCmdOptions{
		name: "tenant-b",
		json: true,
	}

	err := runDescribeNamespaceCmd(context.Background(), svc, options)

	assert.NoError(t, err)
	assert.Equal(t, options.name, svc.lastDescribeArg)
}
