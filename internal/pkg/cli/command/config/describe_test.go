package config

import (
	"errors"
	"testing"

	"github.com/pinecone-io/cli/internal/pkg/cli/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_runDescribeCmd_ReturnsErrorOnUnknownKey(t *testing.T) {
	svc := &mockConfigService{describeErr: errors.New("unknown config key")}

	err := runDescribeCmd(svc, "bad-key", DescribeCmdOptions{})

	assert.Error(t, err)
	assert.Equal(t, "bad-key", svc.lastDescribeKey)
}

func Test_runDescribeCmd_TabularOutput(t *testing.T) {
	svc := &mockConfigService{
		describeResult: ConfigDescription{
			Key:         "environment",
			Value:       "production",
			Description: "Pinecone environment",
			Sensitive:   false,
			ValidValues: []string{"production", "staging"},
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runDescribeCmd(svc, "environment", DescribeCmdOptions{})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "environment")
	assert.Contains(t, out, "production")
}

func Test_runDescribeCmd_JSONOutput(t *testing.T) {
	svc := &mockConfigService{
		describeResult: ConfigDescription{
			Key:         "environment",
			Value:       "production",
			Description: "Pinecone environment",
			Sensitive:   false,
			ValidValues: []string{"production", "staging"},
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runDescribeCmd(svc, "environment", DescribeCmdOptions{json: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"environment"`)
	assert.Contains(t, out, `"production"`)
	assert.Contains(t, out, `"valid_values"`)
}

func Test_runDescribeCmd_MasksSensitiveKeyInJSON(t *testing.T) {
	svc := &mockConfigService{
		describeResult: ConfigDescription{
			Key:       "api-key",
			Value:     "supersecretvalue",
			Sensitive: true,
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runDescribeCmd(svc, "api-key", DescribeCmdOptions{json: true, reveal: false})
		assert.NoError(t, err)
	})

	assert.NotContains(t, out, "supersecretvalue")
}

func Test_runDescribeCmd_RevealsSensitiveKeyInJSON(t *testing.T) {
	svc := &mockConfigService{
		describeResult: ConfigDescription{
			Key:       "api-key",
			Value:     "supersecretvalue",
			Sensitive: true,
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runDescribeCmd(svc, "api-key", DescribeCmdOptions{json: true, reveal: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "supersecretvalue")
}
