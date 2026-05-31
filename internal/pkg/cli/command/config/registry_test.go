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
