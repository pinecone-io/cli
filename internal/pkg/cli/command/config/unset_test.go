package config

import (
	"context"
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_runUnsetCmd_ReturnsErrorOnUnknownKey(t *testing.T) {
	svc := &mockConfigService{unsetErr: errors.New("unknown config key")}

	err := runUnsetCmd(context.Background(), svc, "bad-key", UnsetCmdOptions{})

	assert.Error(t, err)
}

func Test_runUnsetCmd_Succeeds(t *testing.T) {
	svc := &mockConfigService{}

	err := runUnsetCmd(context.Background(), svc, "api-key", UnsetCmdOptions{})

	assert.NoError(t, err)
	assert.Equal(t, "api-key", svc.lastUnsetKey)
}

func Test_runUnsetCmd_SucceedsWithOnChangeLines(t *testing.T) {
	svc := &mockConfigService{
		unsetLines: []string{"You have been logged out"},
	}

	err := runUnsetCmd(context.Background(), svc, "environment", UnsetCmdOptions{})

	assert.NoError(t, err)
	assert.Equal(t, "environment", svc.lastUnsetKey)
}

func Test_runUnsetCmd_JSONOutput(t *testing.T) {
	svc := &mockConfigService{}

	out := testutils.CaptureStdout(t, func() {
		err := runUnsetCmd(context.Background(), svc, "api-key", UnsetCmdOptions{json: true})
		assert.NoError(t, err)
	})

	assert.JSONEq(t, `{"key":"api-key","cleared":true}`, out)
}
