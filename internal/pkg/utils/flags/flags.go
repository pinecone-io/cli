package flags

import (
	"encoding/json"
	"maps"
	"os"
	"strconv"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

type JSONObject map[string]any
type Float32List []float32
type Int32List []int32

func (m *JSONObject) Set(value string) error {
	// allow passing "@file.json" to read a file path and parse as JSON
	if strings.HasPrefix(value, "@") {
		filePath := strings.TrimPrefix(value, "@")
		raw, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		return json.Unmarshal(raw, m)
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
		raw, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		value = string(raw)
	}

	value = strings.TrimSpace(value)
	if value == "" {
		*m = (*m)[:0]
		return nil
	}

	// JSON array
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		var arr []float32
		if err := json.Unmarshal([]byte(value), &arr); err != nil {
			return pcio.Errorf("failed to parse JSON float32 array: %w", err)
		}
		*m = append((*m)[:0], arr...)
		return nil
	}

	// CSV/whitespace separated list
	vals := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
	out := make([]float32, 0, len(vals))
	for _, val := range vals {
		if val == "" {
			continue
		}
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return pcio.Errorf("invalid float32 %q: %w", val, err)
		}
		out = append(out, float32(f))
	}
	*m = append((*m)[:0], out...)
	return nil
}

func (m *Float32List) String() string {
	if m == nil || len(*m) == 0 {
		return ""
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func (*Float32List) Type() string { return "float32 json-array|csv-list" }

func (m *Int32List) Set(value string) error {
	// allow passing "@file.json" to read a file path and parse as JSON
	if strings.HasPrefix(value, "@") {
		filePath := strings.TrimPrefix(value, "@")
		raw, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		return json.Unmarshal(raw, m)
	}

	value = strings.TrimSpace(value)
	if value == "" {
		*m = (*m)[:0]
		return nil
	}

	// JSON array
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		var arr []int32
		if err := json.Unmarshal([]byte(value), &arr); err != nil {
			return pcio.Errorf("failed to parse JSON int32 array: %w", err)
		}
		*m = append((*m)[:0], arr...)
		return nil
	}

	// CSV/whitespace separated list
	vals := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
	out := make([]int32, 0, len(vals))
	for _, val := range vals {
		if val == "" {
			continue
		}
		i, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return pcio.Errorf("invalid int32 %q: %w", val, err)
		}
		out = append(out, int32(i))
	}
	*m = append((*m)[:0], out...)
	return nil
}

func (m *Int32List) String() string {
	if m == nil || len(*m) == 0 {
		return ""
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func (*Int32List) Type() string { return "int32 json-array|csv-list" }
