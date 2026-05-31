package config

import (
	"context"
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_runSetCmd_ReturnsErrorOnUnknownKey(t *testing.T) {
	svc := &mockConfigService{getErr: errors.New("unknown config key")}

	err := runSetCmd(context.Background(), svc, "bad-key", "value", SetCmdOptions{})

	assert.Error(t, err)
	assert.Empty(t, svc.lastSetKey)
}

func Test_runSetCmd_ReturnsNilOnNoChange(t *testing.T) {
	svc := &mockConfigService{
		getValue: "production",
		setErr:   ErrNoChange,
	}

	err := runSetCmd(context.Background(), svc, "environment", "production", SetCmdOptions{})

	assert.NoError(t, err)
	assert.Equal(t, "environment", svc.lastSetKey)
	assert.Equal(t, "production", svc.lastSetValue)
}

func Test_runSetCmd_ReturnsErrorOnValidationFailure(t *testing.T) {
	svc := &mockConfigService{
		getValue: "production",
		setErr:   errors.New("invalid value"),
	}

	err := runSetCmd(context.Background(), svc, "environment", "invalid", SetCmdOptions{})

	assert.Error(t, err)
	assert.Equal(t, "environment", svc.lastSetKey)
}

func Test_runSetCmd_Succeeds(t *testing.T) {
	svc := &mockConfigService{getValue: "production"}

	err := runSetCmd(context.Background(), svc, "environment", "staging", SetCmdOptions{})

	assert.NoError(t, err)
	assert.Equal(t, "environment", svc.lastSetKey)
	assert.Equal(t, "staging", svc.lastSetValue)
}

func Test_runSetCmd_SucceedsWithOnChangeLines(t *testing.T) {
	svc := &mockConfigService{
		getValue: "production",
		setLines: []string{"You have been logged out", "API key cleared"},
	}

	err := runSetCmd(context.Background(), svc, "environment", "staging", SetCmdOptions{})

	assert.NoError(t, err)
	assert.Equal(t, "staging", svc.lastSetValue)
}

func Test_runSetCmd_JSONOutput(t *testing.T) {
	svc := &mockConfigService{getValue: "production"}

	out := testutils.CaptureStdout(t, func() {
		err := runSetCmd(context.Background(), svc, "environment", "staging", SetCmdOptions{json: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"environment"`)
	assert.Contains(t, out, `"staging"`)
}

func Test_runSetCmd_JSONOutputOnNoChange(t *testing.T) {
	svc := &mockConfigService{
		getValue: "production",
		setErr:   ErrNoChange,
	}

	out := testutils.CaptureStdout(t, func() {
		err := runSetCmd(context.Background(), svc, "environment", "production", SetCmdOptions{json: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"environment"`)
	assert.Contains(t, out, `"production"`)
}
