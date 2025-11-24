package inputpolicy

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DefaultMaxBodyJSONBytes int64 = 1 << 30 // 1 GiB default cap
)

var (
	// Limits can be overridden via env vars at process start.
	MaxBodyJSONBytes int64 = parseSizeFromEnv("PC_CLI_MAX_JSON_BYTES", DefaultMaxBodyJSONBytes)
)

func parseSizeFromEnv(key string, def int64) int64 {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def
	}
	// support simple suffixes: k, m, g (base 1024)
	mult := int64(1)
	switch last := strings.ToLower(val[len(val)-1:]); last {
	case "k":
		mult = 1 << 10
		val = val[:len(val)-1]
	case "m":
		mult = 1 << 20
		val = val[:len(val)-1]
	case "g":
		mult = 1 << 30
		val = val[:len(val)-1]
	}
	n, err := strconv.ParseInt(val, 10, 64)
	if err != nil || n <= 0 {
		return def
	}
	return n * mult
}

// Minimal path validation: must exist, not be a directory, and be a regular file.
// Symlinks are resolved before checking.
func ValidatePath(path string) error {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path %s: %w", path, err)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return fmt.Errorf("failed to stat path %s: %w", resolved, err)
	}
	if info.IsDir() {
		return fmt.Errorf("path %s is a directory", resolved)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("path %s is not a regular file", resolved)
	}
	return nil
}
