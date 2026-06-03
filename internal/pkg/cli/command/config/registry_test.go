package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupKey_ValidKeys(t *testing.T) {
	for _, key := range configKeys {
		t.Run(key, func(t *testing.T) {
			desc, err := lookupKey(key)
			assert.NoError(t, err)
			assert.NotEmpty(t, desc.Description)
		})
	}
}

func TestLookupKey_InvalidKey(t *testing.T) {
	_, err := lookupKey("not-a-real-key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not-a-real-key")
	for _, key := range configKeys {
		assert.Contains(t, err.Error(), key)
	}
}

func TestConfigKeysMatchRegistry(t *testing.T) {
	for _, key := range configKeys {
		_, err := lookupKey(key)
		assert.NoError(t, err, "key %q is in configKeys but not in configRegistry", key)
	}
	assert.Equal(t, len(configKeys), len(configRegistry),
		"configRegistry has keys not listed in configKeys")
}

func TestVisibleKeysFiltersHiddenKeys(t *testing.T) {
	visibleKeys := visibleKeys()
	for key, desc := range configRegistry {
		if desc.Hidden {
			assert.NotContains(t, visibleKeys, key)
		} else {
			assert.Contains(t, visibleKeys, key)
		}
	}
}

func TestDisplayValue_Empty(t *testing.T) {
	assert.Equal(t, "<not set>", displayValue(""))
}

func TestDisplayValue_NonEmpty(t *testing.T) {
	assert.Equal(t, "production", displayValue("production"))
}

// Service-level env var override tests.
//
// Get, List, and Describe are read-only and safe to run without any file
// setup — GetStored returns the key's default value when no config file
// exists. Set and Unset call persistStr which writes to the real config
// file, so their env-var-skip path is not unit-tested here; it is covered
// by the e2e suite.

func TestDefaultConfigService_Get_ReturnsEnvVarValueAndSetsOverride(t *testing.T) {
	t.Setenv("PINECONE_ENVIRONMENT", "staging")
	svc := newDefaultConfigService()

	value, _, envVarName, envVarOverride, err := svc.Get("environment")

	assert.NoError(t, err)
	assert.Equal(t, "staging", value)
	assert.Equal(t, "PINECONE_ENVIRONMENT", envVarName)
	assert.True(t, envVarOverride)
}

func TestDefaultConfigService_Get_ReportsNoOverrideWhenEnvVarAbsent(t *testing.T) {
	t.Setenv("PINECONE_ENVIRONMENT", "")
	svc := newDefaultConfigService()

	_, _, envVarName, envVarOverride, err := svc.Get("environment")

	assert.NoError(t, err)
	assert.Equal(t, "PINECONE_ENVIRONMENT", envVarName) // always present for bound keys
	assert.False(t, envVarOverride)
}

func TestDefaultConfigService_Get_EnvVarNameEmptyForUnboundKey(t *testing.T) {
	svc := newDefaultConfigService()

	_, _, envVarName, envVarOverride, err := svc.Get("color")

	assert.NoError(t, err)
	assert.Empty(t, envVarName)
	assert.False(t, envVarOverride)
}

func TestDefaultConfigService_List_ReflectsActiveEnvVarOverride(t *testing.T) {
	t.Setenv("PINECONE_ENVIRONMENT", "staging")
	svc := newDefaultConfigService()

	entries := svc.List(true) // include hidden keys

	for _, e := range entries {
		if e.Key == "environment" {
			assert.Equal(t, "staging", e.Value)
			assert.Equal(t, "PINECONE_ENVIRONMENT", e.EnvVarName)
			assert.True(t, e.EnvVarOverride)
			return
		}
	}
	t.Fatal("environment entry not found in list")
}

func TestDefaultConfigService_List_ReportsNoOverrideWhenEnvVarAbsent(t *testing.T) {
	t.Setenv("PINECONE_ENVIRONMENT", "")
	svc := newDefaultConfigService()

	entries := svc.List(true)

	for _, e := range entries {
		if e.Key == "environment" {
			assert.Equal(t, "PINECONE_ENVIRONMENT", e.EnvVarName)
			assert.False(t, e.EnvVarOverride)
			return
		}
	}
	t.Fatal("environment entry not found in list")
}

func TestDefaultConfigService_List_EnvVarNameEmptyForUnboundKey(t *testing.T) {
	svc := newDefaultConfigService()

	for _, e := range svc.List(false) {
		if e.Key == "color" {
			assert.Empty(t, e.EnvVarName)
			assert.False(t, e.EnvVarOverride)
			return
		}
	}
	t.Fatal("color entry not found in list")
}

func TestDefaultConfigService_Describe_ReflectsActiveEnvVarOverride(t *testing.T) {
	t.Setenv("PINECONE_API_KEY", "test-api-key")
	svc := newDefaultConfigService()

	desc, err := svc.Describe("api-key")

	assert.NoError(t, err)
	assert.Equal(t, "test-api-key", desc.Value)
	assert.Equal(t, "PINECONE_API_KEY", desc.EnvVarName)
	assert.True(t, desc.EnvVarOverride)
}

func TestDefaultConfigService_Describe_ReportsNoOverrideWhenEnvVarAbsent(t *testing.T) {
	t.Setenv("PINECONE_API_KEY", "")
	svc := newDefaultConfigService()

	desc, err := svc.Describe("api-key")

	assert.NoError(t, err)
	assert.Equal(t, "PINECONE_API_KEY", desc.EnvVarName)
	assert.False(t, desc.EnvVarOverride)
}

func TestDefaultConfigService_Describe_EnvVarNameEmptyForUnboundKey(t *testing.T) {
	svc := newDefaultConfigService()

	desc, err := svc.Describe("color")

	assert.NoError(t, err)
	assert.Empty(t, desc.EnvVarName)
	assert.False(t, desc.EnvVarOverride)
}
