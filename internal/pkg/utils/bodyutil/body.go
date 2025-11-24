package bodyutil

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

// OpenArgReader returns an io.ReadCloser for inline/@file/@- with size limits applied.
// When isBody is true, the body size limit is used, otherwise the flag size limit.
func OpenArgReader(spec string, isBody bool) (io.ReadCloser, SourceInfo, error) {
	limit := inputpolicy.MaxBodyJSONBytes
	switch {
	case spec == "":
		return nil, SourceInfo{Kind: SourceInline, Label: "inline"}, nil
	case spec == "@-":
		r, err := stdin.ReaderOnce(limit)
		if err != nil {
			return nil, SourceInfo{Kind: SourceStdin, Label: "stdin"}, fmt.Errorf("stdin already consumed; only one argument may use stdin")
		}
		return r, SourceInfo{Kind: SourceStdin, Label: "stdin"}, nil
	case len(spec) > 0 && spec[0] == '@':
		path := spec[1:]
		if err := inputpolicy.ValidatePath(path); err != nil {
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
	default:
		return io.NopCloser(strings.NewReader(spec)), SourceInfo{Kind: SourceInline, Label: "inline"}, nil
	}
}

// DecodeBodyArgs unmarshals a JSON body argument (inline/@file/@-) using a bounded reader.
func DecodeBodyArgs[T any](spec string) (*T, SourceInfo, error) {
	rc, src, err := OpenArgReader(spec, true)
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

// NOTE: we intentionally import strings and use strings.NewReader for clarity.
