package flags

import (
	"encoding/json"
	"io"
	"maps"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/inputpolicy"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/stdin"
)

type JSONObject map[string]any
type Float32List []float32
type UInt32List []uint32
type StringList []string

func (m *JSONObject) Set(value string) error {
	// allow passing "@file.json" to read a file path and parse as JSON
	if strings.HasPrefix(value, "@") {
		filePath := strings.TrimPrefix(value, "@")
		if filePath == "-" {
			rc, err := stdin.ReaderOnce(inputpolicy.MaxBodyJSONBytes)
			if err != nil {
				if err == io.ErrUnexpectedEOF {
					return pcio.Errorf("stdin already consumed; only one argument may use stdin")
				}
				return err
			}
			defer rc.Close()
			dec := json.NewDecoder(rc)
			return dec.Decode(m)
		} else {
			if err := inputpolicy.ValidatePath(filePath); err != nil {
				return err
			}
			f, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()
			dec := json.NewDecoder(io.LimitReader(f, inputpolicy.MaxBodyJSONBytes))
			return dec.Decode(m)
		}
	}

	var tmp map[string]any
	if err := json.Unmarshal([]byte(value), &tmp); err != nil {
		return pcio.Errorf("failed to parse JSON: %w", err)
	}
	if *m == nil {
		*m = make(map[string]any)
	}
	maps.Copy((*m), tmp)
	return nil
}

func (m *JSONObject) String() string {
	if m == nil || len(*m) == 0 {
		return ""
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func (*JSONObject) Type() string { return "json-object" }

func (m *Float32List) Set(value string) error {
	// allow passing "@file.json" to read a file path and parse as JSON
	if strings.HasPrefix(value, "@") {
		filePath := strings.TrimPrefix(value, "@")
		if filePath == "-" {
			rc, err := stdin.ReaderOnce(inputpolicy.MaxBodyJSONBytes)
			if err != nil {
				if err == io.ErrUnexpectedEOF {
					return pcio.Errorf("stdin already consumed; only one argument may use stdin")
				}
				return err
			}
			defer rc.Close()
			var arr []float32
			if err := json.NewDecoder(rc).Decode(&arr); err != nil {
				return pcio.Errorf("failed to parse JSON float32 array: %w", err)
			}
			*m = append((*m)[:0], arr...)
			return nil
		} else {
			if err := inputpolicy.ValidatePath(filePath); err != nil {
				return err
			}
			f, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()
			var arr []float32
			if err := json.NewDecoder(io.LimitReader(f, inputpolicy.MaxBodyJSONBytes)).Decode(&arr); err != nil {
				return pcio.Errorf("failed to parse JSON float32 array: %w", err)
			}
			*m = append((*m)[:0], arr...)
			return nil
		}
	}

	value = strings.TrimSpace(value)
	if value == "" {
		*m = (*m)[:0]
		return nil
	}

	var arr []float32
	if err := json.Unmarshal([]byte(value), &arr); err != nil {
		return pcio.Errorf("failed to parse JSON float32 array: %w", err)
	}
	*m = append((*m)[:0], arr...)
	return nil
}

func (m *Float32List) String() string {
	if m == nil || len(*m) == 0 {
		return ""
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func (*Float32List) Type() string { return "float32 json-array" }

func (m *UInt32List) Set(value string) error {
	// allow passing "@file.json" to read a file path and parse as JSON
	if strings.HasPrefix(value, "@") {
		filePath := strings.TrimPrefix(value, "@")
		if filePath == "-" {
			rc, err := stdin.ReaderOnce(inputpolicy.MaxBodyJSONBytes)
			if err != nil {
				if err == io.ErrUnexpectedEOF {
					return pcio.Errorf("stdin already consumed; only one argument may use stdin")
				}
				return err
			}
			defer rc.Close()
			var arr []uint32
			if err := json.NewDecoder(rc).Decode(&arr); err != nil {
				return pcio.Errorf("failed to parse JSON uint32 array: %w", err)
			}
			*m = append((*m)[:0], arr...)
			return nil
		} else {
			if err := inputpolicy.ValidatePath(filePath); err != nil {
				return err
			}
			f, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()
			var arr []uint32
			if err := json.NewDecoder(io.LimitReader(f, inputpolicy.MaxBodyJSONBytes)).Decode(&arr); err != nil {
				return pcio.Errorf("failed to parse JSON uint32 array: %w", err)
			}
			*m = append((*m)[:0], arr...)
			return nil
		}
	}

	value = strings.TrimSpace(value)
	if value == "" {
		*m = (*m)[:0]
		return nil
	}

	var arr []uint32
	if err := json.Unmarshal([]byte(value), &arr); err != nil {
		return pcio.Errorf("failed to parse JSON uint32 array: %w", err)
	}
	*m = append((*m)[:0], arr...)
	return nil
}

func (m *UInt32List) String() string {
	if m == nil || len(*m) == 0 {
		return ""
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func (*UInt32List) Type() string { return "uint32 json-array" }

func (m *StringList) Set(value string) error {
	// allow passing "@file.json" to read a file path and parse as JSON
	if strings.HasPrefix(value, "@") {
		filePath := strings.TrimPrefix(value, "@")
		if filePath == "-" {
			rc, err := stdin.ReaderOnce(inputpolicy.MaxBodyJSONBytes)
			if err != nil {
				if err == io.ErrUnexpectedEOF {
					return pcio.Errorf("stdin already consumed; only one argument may use stdin")
				}
				return err
			}
			defer rc.Close()
			var arr []string
			if err := json.NewDecoder(rc).Decode(&arr); err != nil {
				return pcio.Errorf("failed to parse JSON string array: %w", err)
			}
			*m = append((*m)[:0], arr...)
			return nil
		} else {
			if err := inputpolicy.ValidatePath(filePath); err != nil {
				return err
			}
			f, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()
			var arr []string
			if err := json.NewDecoder(io.LimitReader(f, inputpolicy.MaxBodyJSONBytes)).Decode(&arr); err != nil {
				return pcio.Errorf("failed to parse JSON string array: %w", err)
			}
			*m = append((*m)[:0], arr...)
			return nil
		}
	}
	value = strings.TrimSpace(value)
	if value == "" {
		*m = (*m)[:0]
		return nil
	}
	var arr []string
	if err := json.Unmarshal([]byte(value), &arr); err != nil {
		return pcio.Errorf("failed to parse JSON string array: %w", err)
	}
	*m = append((*m)[:0], arr...)
	return nil
}

func (m *StringList) String() string {
	if m == nil || len(*m) == 0 {
		return ""
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func (*StringList) Type() string { return "string json-array" }
