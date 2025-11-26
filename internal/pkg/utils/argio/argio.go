package argio

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/inputpolicy"
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

// OpenReader returns an io.ReadCloser for inline text, JSON/JSONL files, or stdin.
// Supported syntax:
//   - Inline JSON (default)
//   - File paths ending with .json or .jsonl
//   - "-" for stdin (only one consumer allowed)
func OpenReader(value string) (io.ReadCloser, SourceInfo, error) {
	limit := inputpolicy.MaxBodyJSONBytes
	switch {
	case value == "": // empty value is inline
		return io.NopCloser(strings.NewReader("")), SourceInfo{Kind: SourceInline, Label: "inline"}, nil
	case value == "-": // stdin
		r, err := stdin.ReaderOnce(limit)
		if err != nil {
			return nil, SourceInfo{Kind: SourceStdin, Label: "stdin"}, fmt.Errorf("stdin already consumed; only one argument may use '-' per command")
		}

		return r, SourceInfo{Kind: SourceStdin, Label: "stdin"}, nil
	case looksLikeJSONFile(value):
		return openJSONFile(value, limit)
	default: // if no stdin and no file, it's inline
		return io.NopCloser(strings.NewReader(value)), SourceInfo{Kind: SourceInline, Label: "inline"}, nil
	}
}

// ReadAll reads the entire argument from the inline/file/stdin value into memory using a bounded reader.
func ReadAll(value string) ([]byte, SourceInfo, error) {
	rc, src, err := OpenReader(value)
	if err != nil {
		return nil, src, err
	}
	defer rc.Close()

	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, src, err
	}

	return b, src, nil
}

// DecodeJSONArg unmarshals a JSON argument from the inline/file/stdin value into a generic type using a bounded reader.
func DecodeJSONArg[T any](value string) (*T, SourceInfo, error) {
	rc, src, err := OpenReader(value)
	if err != nil {
		return nil, src, err
	}

	var closer io.Closer
	if rc != nil {
		closer = rc
		defer closer.Close()
	}

	dec := json.NewDecoder(rc)
	dec.DisallowUnknownFields()
	var out T
	if err := dec.Decode(&out); err != nil {
		return nil, src, fmt.Errorf("invalid JSON from %s: %w", src.Label, err)
	}

	return &out, src, nil
}

func looksLikeJSONFile(value string) bool {
	lower := strings.ToLower(value)
	return strings.HasSuffix(lower, ".json") || strings.HasSuffix(lower, ".jsonl")
}

func openJSONFile(path string, limit int64) (io.ReadCloser, SourceInfo, error) {
	if err := inputpolicy.ValidatePath(path); err != nil {
		return nil, SourceInfo{Kind: SourceFile, Label: path}, err
	}

	if _, err := os.Stat(path); err != nil {
		return nil, SourceInfo{Kind: SourceFile, Label: path}, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, SourceInfo{Kind: SourceFile, Label: path}, err
	}

	return struct {
		io.Reader
		io.Closer
	}{Reader: io.LimitReader(f, limit), Closer: f}, SourceInfo{Kind: SourceFile, Label: path}, nil
}
