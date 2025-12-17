package namespace

import (
	"context"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/stretchr/testify/assert"
)

func Test_runCreateNamespaceCmd_RequiresName(t *testing.T) {
	svc := &mockNamespaceService{}
	options := createNamespaceCmdOptions{
		name: "",
	}

	err := runCreateNamespaceCmd(context.Background(), svc, options)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--name is required")
	assert.Nil(t, svc.lastCreateReq)
}

func Test_runCreateNamespaceCmd_Succeeds(t *testing.T) {
	svc := &mockNamespaceService{
		createResp: &mockNamespaceDescription,
	}
	options := createNamespaceCmdOptions{
		name:           "tenant-a",
		metadataSchema: []string{"category:keyword", "brand:keyword"},
	}

	err := runCreateNamespaceCmd(context.Background(), svc, options)

	if assert.NoError(t, err) {
		if assert.NotNil(t, svc.lastCreateReq) {
			assert.Equal(t, options.name, svc.lastCreateReq.Name)
			expectedSchema := sdk.BuildMetadataSchema(options.metadataSchema)
			assert.Equal(t, expectedSchema, svc.lastCreateReq.Schema)
		}
	}
}

func Test_runCreateNamespaceCmd_SucceedsJSON(t *testing.T) {
	svc := &mockNamespaceService{
		createResp: &mockNamespaceDescription,
	}
	options := createNamespaceCmdOptions{
		name: "tenant-b",
		json: true,
	}

	err := runCreateNamespaceCmd(context.Background(), svc, options)

	if assert.NoError(t, err) {
		if assert.NotNil(t, svc.lastCreateReq) {
			assert.Equal(t, options.name, svc.lastCreateReq.Name)
		}
	}
}
