package namespace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_runListNamespaceCmd_RequiresIndexName(t *testing.T) {
	svc := &mockNamespaceService{}
	options := listNamespaceCmdOptions{
		indexName: "",
	}

	err := runListNamespaceCmd(context.Background(), svc, options)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--index-name is required")
	assert.Nil(t, svc.lastListParams)
}

func Test_runListNamespaceCmd_Succeeds(t *testing.T) {
	svc := &mockNamespaceService{
		listResp: &mockListNamespacesResponse,
	}
	options := listNamespaceCmdOptions{
		indexName:       "my-index",
		paginationToken: "token-1",
		limit:           10,
		prefix:          "tenant-",
	}

	err := runListNamespaceCmd(context.Background(), svc, options)

	if assert.NoError(t, err) {
		if assert.NotNil(t, svc.lastListParams) {
			if assert.NotNil(t, svc.lastListParams.Limit) {
				assert.Equal(t, uint32(10), *svc.lastListParams.Limit)
			}
			if assert.NotNil(t, svc.lastListParams.PaginationToken) {
				assert.Equal(t, "token-1", *svc.lastListParams.PaginationToken)
			}
			if assert.NotNil(t, svc.lastListParams.Prefix) {
				assert.Equal(t, "tenant-", *svc.lastListParams.Prefix)
			}
		}
	}
}

func Test_runListNamespaceCmd_SucceedsJSON(t *testing.T) {
	svc := &mockNamespaceService{
		listResp: &mockListNamespacesResponse,
	}
	options := listNamespaceCmdOptions{
		indexName: "my-index",
		json:      true,
	}

	err := runListNamespaceCmd(context.Background(), svc, options)

	assert.NoError(t, err)
	assert.NotNil(t, svc.lastListParams)
	if assert.NotNil(t, svc.lastListParams) {
		assert.Nil(t, svc.lastListParams.Limit)
		assert.Nil(t, svc.lastListParams.PaginationToken)
		assert.Nil(t, svc.lastListParams.Prefix)
	}
}
