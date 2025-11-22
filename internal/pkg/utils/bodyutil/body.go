package bodyutil

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/stdin"
)

type ArgSource int

const (
	SourceInline ArgSource = iota
	SourceFile
	SourceStdin
)

type SourceInfo struct {
	Kind  ArgSource `json:"kind"`  // SourceInline, SourceFile, SourceStdin
	Label string    `json:"label"` // "inline", file path, "stdin"
}

// DecodeBodyArgs unmarshals a JSON body argument (inline/@file/@-) into the provided generic type.
func DecodeBodyArgs[T any](spec string) (*T, SourceInfo, error) {
	b, src, err := ReadArg(spec)
	if err != nil || len(b) == 0 {
		return nil, src, err
	}
	var out T
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, src, fmt.Errorf("invalid JSON from %s: %w", src.Label, err)
	}
	return &out, src, nil
}

// ReadArg reads inline JSON, @file, or @- (stdin once) and returns bytes and a source label.
func ReadArg(spec string) ([]byte, SourceInfo, error) {
	switch {
	case spec == "":
		return nil, SourceInfo{Kind: SourceInline, Label: "inline"}, nil
	case spec == "@-":
		b, err := stdin.ReadAllOnce()
		if err != nil {
			return nil, SourceInfo{Kind: SourceStdin, Label: "stdin"}, fmt.Errorf("stdin already consumed; only one argument may use stdin")
		}
		return b, SourceInfo{Kind: SourceStdin, Label: "stdin"}, nil
	case len(spec) > 0 && spec[0] == '@':
		path := spec[1:]
		b, err := os.ReadFile(path)
		return b, SourceInfo{Kind: SourceFile, Label: path}, err
	default:
		return []byte(spec), SourceInfo{Kind: SourceInline, Label: "inline"}, nil
	}
}
