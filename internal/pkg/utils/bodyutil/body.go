package bodyutil

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/stdin"
)

// DecodeBodyArgs unmarshals a JSON body argument (inline/@file/@-) into the provided generic type.
func DecodeBodyArgs[T any](spec string) (*T, string, error) {
	b, src, err := ReadArg(spec)
	if err != nil || len(b) == 0 {
		return nil, src, err
	}
	var out T
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, src, fmt.Errorf("invalid JSON from %s: %w", src, err)
	}
	return &out, src, nil
}

// ReadArg reads inline JSON, @file, or @- (stdin once) and returns bytes and a source label.
func ReadArg(spec string) ([]byte, string, error) {
	switch {
	case spec == "":
		return nil, "", nil
	case spec == "@-":
		b, err := stdin.ReadAllOnce()
		if err != nil {
			return nil, "stdin", fmt.Errorf("stdin already consumed; only one argument may use stdin")
		}
		return b, "stdin", nil
	case len(spec) > 0 && spec[0] == '@':
		path := spec[1:]
		b, err := os.ReadFile(path)
		return b, path, err
	default:
		return []byte(spec), "inline", nil
	}
}
