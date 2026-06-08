package importcmd

import (
	"context"
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/stretchr/testify/assert"
)

func Test_runListImportsCmd_PopulatesParams(t *testing.T) {
	svc := &mockImportService{
		listImportsResp: &pinecone.ListImportsResponse{},
	}
	opts := listImportsCmdOptions{
		limit:           5,
		paginationToken: "next-page",
	}

	err := runListImportsCmd(context.Background(), svc, opts)

	if assert.NoError(t, err) {
		if assert.NotNil(t, svc.lastListImportsLimit) {
			assert.Equal(t, int32(5), *svc.lastListImportsLimit)
		}
		if assert.NotNil(t, svc.lastListImportsPaginationToken) {
			assert.Equal(t, "next-page", *svc.lastListImportsPaginationToken)
		}
	}
}

func Test_runListImportsCmd_OmitsZeroLimit(t *testing.T) {
	svc := &mockImportService{
		listImportsResp: &pinecone.ListImportsResponse{},
	}
	opts := listImportsCmdOptions{limit: 0}

	err := runListImportsCmd(context.Background(), svc, opts)

	if assert.NoError(t, err) {
		assert.Nil(t, svc.lastListImportsLimit)
	}
}

func Test_runListImportsCmd_SucceedsJSON(t *testing.T) {
	svc := &mockImportService{
		listImportsResp: &pinecone.ListImportsResponse{
			Imports: []*pinecone.Import{{Id: "import-1"}},
		},
	}
	opts := listImportsCmdOptions{json: true}

	out := testutils.CaptureStdout(t, func() {
		err := runListImportsCmd(context.Background(), svc, opts)
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"import-1"`)
}

func Test_runListImportsCmd_PropagatesError(t *testing.T) {
	svc := &mockImportService{
		listImportsErr: errors.New("list failed"),
	}
	opts := listImportsCmdOptions{}

	err := runListImportsCmd(context.Background(), svc, opts)

	assert.EqualError(t, err, "list failed")
}
