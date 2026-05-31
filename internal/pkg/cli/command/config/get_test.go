package config

import (
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_runGetCmd_ReturnsErrorOnUnknownKey(t *testing.T) {
	svc := &mockConfigService{getErr: errors.New("unknown config key")}

	err := runGetCmd(svc, "bad-key", GetCmdOptions{})

	assert.Error(t, err)
	assert.Equal(t, "bad-key", svc.lastGetKey)
}

func Test_runGetCmd_Succeeds(t *testing.T) {
	svc := &mockConfigService{getValue: "production"}

	err := runGetCmd(svc, "environment", GetCmdOptions{})

	assert.NoError(t, err)
	assert.Equal(t, "environment", svc.lastGetKey)
}

func Test_runGetCmd_JSONOutput(t *testing.T) {
	svc := &mockConfigService{getValue: "production"}

	out := testutils.CaptureStdout(t, func() {
		err := runGetCmd(svc, "environment", GetCmdOptions{json: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"environment"`)
	assert.Contains(t, out, `"production"`)
}

func Test_runGetCmd_MasksSensitiveKeyInJSON(t *testing.T) {
	svc := &mockConfigService{getValue: "supersecretvalue", getSensitive: true}

	out := testutils.CaptureStdout(t, func() {
		err := runGetCmd(svc, "api-key", GetCmdOptions{json: true, reveal: false})
		assert.NoError(t, err)
	})

	assert.NotContains(t, out, "supersecretvalue")
}

func Test_runGetCmd_RevealsSensitiveKeyInJSON(t *testing.T) {
	svc := &mockConfigService{getValue: "supersecretvalue", getSensitive: true}

	out := testutils.CaptureStdout(t, func() {
		err := runGetCmd(svc, "api-key", GetCmdOptions{json: true, reveal: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "supersecretvalue")
}
