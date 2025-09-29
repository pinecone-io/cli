package cache

import (
	"path/filepath"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
)

var (
	// Initialize cache in the config directory
	Cache = NewFileCache(filepath.Join(configuration.ConfigDirPath(), "cache"))
)

// GetOrFetch is a helper function to get cached data or fetch from API
func GetOrFetch[T any](key string, ttl time.Duration, fetchFunc func() (*T, error)) (*T, error) {
	var cached T
	if found, err := Cache.Get(key, &cached); found && err == nil {
		return &cached, nil
	}

	// Fetch from API
	data, err := fetchFunc()
	if err != nil {
		return nil, err
	}

	// Cache the result
	Cache.Set(key, data, ttl)
	return data, nil
}

// CacheWithTTL is a helper function to cache data with a specific TTL
func CacheWithTTL(key string, data interface{}, ttl time.Duration) error {
	return Cache.Set(key, data, ttl)
}

// GetCached is a helper function to get cached data
func GetCached[T any](key string) (*T, bool, error) {
	var cached T
	found, err := Cache.Get(key, &cached)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}
	return &cached, true, nil
}
