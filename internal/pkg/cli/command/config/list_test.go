package config

import (
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_runListCmd_TabularOutputIncludesHeader(t *testing.T) {
	svc := &mockConfigService{listResult: []ConfigEntry{}}

	out := testutils.CaptureStdout(t, func() {
		err := runListCmd(svc, ListCmdOptions{})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "KEY")
	assert.Contains(t, out, "VALUE")
	assert.Contains(t, out, "DESCRIPTION")
}

func Test_runListCmd_TabularOutputMasksSensitiveKey(t *testing.T) {
	svc := &mockConfigService{
		listResult: []ConfigEntry{
			{Key: "api-key", Value: "sk-supersecret", Description: "API key", Sensitive: true},
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runListCmd(svc, ListCmdOptions{})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "api-key")
	assert.NotContains(t, out, "sk-supersecret")
}

func Test_runListCmd_TabularOutputRevealsSensitiveKey(t *testing.T) {
	svc := &mockConfigService{
		listResult: []ConfigEntry{
			{Key: "api-key", Value: "sk-supersecret", Sensitive: true},
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runListCmd(svc, ListCmdOptions{reveal: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "sk-supersecret")
}

func Test_runListCmd_JSONOutput(t *testing.T) {
	svc := &mockConfigService{
		listResult: []ConfigEntry{
			{Key: "api-key", Value: "sk-supersecret", Description: "API key", Sensitive: true},
			{Key: "color", Value: "true", Description: "Color output", Sensitive: false},
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runListCmd(svc, ListCmdOptions{json: true})
		assert.NoError(t, err)
	})

	// Sensitive key should be masked in JSON output
	assert.NotContains(t, out, "sk-supersecret")
	// Non-sensitive values should appear
	assert.Contains(t, out, `"color"`)
	assert.Contains(t, out, `"true"`)
}

func Test_runListCmd_AllFlagIncludesHiddenKeys(t *testing.T) {
	svc := &mockConfigService{
		listResult: []ConfigEntry{
			{Key: "api-key", Value: "", Description: "API key", Sensitive: true},
			{Key: "color", Value: "true", Description: "Color output"},
			{Key: "environment", Value: "production", Description: "Environment", Hidden: true},
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runListCmd(svc, ListCmdOptions{all: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "environment")
}

func Test_runListCmd_JSONAllFlagIncludesHiddenField(t *testing.T) {
	svc := &mockConfigService{
		listResult: []ConfigEntry{
			{Key: "color", Value: "true", Description: "Color output", Hidden: false},
			{Key: "environment", Value: "production", Description: "Environment", Hidden: true},
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runListCmd(svc, ListCmdOptions{json: true, all: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "environment")
	assert.Contains(t, out, `"hidden": true`)
	// Non-hidden keys should not have the hidden field (omitempty)
	assert.NotContains(t, out, `"hidden": false`)
}

func Test_runListCmd_JSONOutputRevealsSensitiveKey(t *testing.T) {
	svc := &mockConfigService{
		listResult: []ConfigEntry{
			{Key: "api-key", Value: "sk-supersecret", Sensitive: true},
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runListCmd(svc, ListCmdOptions{json: true, reveal: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "sk-supersecret")
}
