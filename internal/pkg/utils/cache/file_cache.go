package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type FileCache struct {
	basePath string
}

type CacheEntry struct {
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
	TTL       time.Duration   `json:"ttl"`
}

func NewFileCache(basePath string) *FileCache {
	return &FileCache{
		basePath: basePath,
	}
}

func (fc *FileCache) Get(key string, target interface{}) (bool, error) {
	cacheFile := filepath.Join(fc.basePath, key+".json")

	// Check if file exists
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return false, nil
	}

	// Read and parse cache file
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return false, err
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return false, err
	}

	// Check if expired
	if time.Since(entry.Timestamp) > entry.TTL {
		os.Remove(cacheFile) // Clean up expired file
		return false, nil
	}

	// Unmarshal data directly into target
	if err := json.Unmarshal(entry.Data, target); err != nil {
		return false, err
	}

	return true, nil
}

func (fc *FileCache) Set(key string, data interface{}, ttl time.Duration) error {
	// Ensure cache directory exists
	if err := os.MkdirAll(fc.basePath, 0755); err != nil {
		return err
	}

	// Marshal the data to JSON first
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	entry := CacheEntry{
		Data:      json.RawMessage(dataBytes),
		Timestamp: time.Now(),
		TTL:       ttl,
	}

	entryBytes, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	cacheFile := filepath.Join(fc.basePath, key+".json")
	return os.WriteFile(cacheFile, entryBytes, 0644)
}

func (fc *FileCache) Delete(key string) error {
	cacheFile := filepath.Join(fc.basePath, key+".json")
	return os.Remove(cacheFile)
}

func (fc *FileCache) Clear() error {
	return os.RemoveAll(fc.basePath)
}
