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

func Test_runDescribeCmd_TabularOutputShowsEnvVarRows(t *testing.T) {
	svc := &mockConfigService{
		describeResult: ConfigDescription{
			Key:            "environment",
			Value:          "staging",
			EnvVarName:     "PINECONE_ENVIRONMENT",
			EnvVarOverride: true,
			Description:    "Pinecone environment",
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runDescribeCmd(svc, "environment", DescribeCmdOptions{})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "ENV VAR NAME")
	assert.Contains(t, out, "$PINECONE_ENVIRONMENT")
	assert.Contains(t, out, "ENV VAR OVERRIDE")
	assert.Contains(t, out, "true")
}

func Test_runDescribeCmd_TabularOutputOmitsEnvVarRowsWhenUnbound(t *testing.T) {
	svc := &mockConfigService{
		describeResult: ConfigDescription{
			Key:         "color",
			Value:       "true",
			Description: "Enable or disable colored terminal output",
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runDescribeCmd(svc, "color", DescribeCmdOptions{})
		assert.NoError(t, err)
	})

	assert.NotContains(t, out, "ENV VAR NAME")
	assert.NotContains(t, out, "ENV VAR OVERRIDE")
}

func Test_runDescribeCmd_JSONOutputIncludesEnvVarFields(t *testing.T) {
	svc := &mockConfigService{
		describeResult: ConfigDescription{
			Key:            "environment",
			Value:          "staging",
			EnvVarName:     "PINECONE_ENVIRONMENT",
			EnvVarOverride: true,
			Description:    "Pinecone environment",
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runDescribeCmd(svc, "environment", DescribeCmdOptions{json: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"PINECONE_ENVIRONMENT"`)
	assert.Contains(t, out, `"env_var_override": true`)
}

func Test_runDescribeCmd_JSONOutputEnvVarOverrideIsFalseWhenNotActive(t *testing.T) {
	svc := &mockConfigService{
		describeResult: ConfigDescription{
			Key:            "environment",
			Value:          "production",
			EnvVarName:     "PINECONE_ENVIRONMENT",
			EnvVarOverride: false,
			Description:    "Pinecone environment",
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runDescribeCmd(svc, "environment", DescribeCmdOptions{json: true})
		assert.NoError(t, err)
	})

	assert.Contains(t, out, `"PINECONE_ENVIRONMENT"`)
	assert.Contains(t, out, `"env_var_override": false`)
}

func Test_runDescribeCmd_JSONOutputOmitsEnvVarFieldsWhenUnbound(t *testing.T) {
	svc := &mockConfigService{
		describeResult: ConfigDescription{
			Key:         "color",
			Value:       "true",
			Description: "Enable or disable colored terminal output",
		},
	}

	out := testutils.CaptureStdout(t, func() {
		err := runDescribeCmd(svc, "color", DescribeCmdOptions{json: true})
		assert.NoError(t, err)
	})

	assert.NotContains(t, out, "env_var_name")
	assert.NotContains(t, out, "env_var_override")
}
