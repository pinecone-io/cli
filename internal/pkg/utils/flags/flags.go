package flags

import (
	"encoding/json"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/argio"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

type JSONObject map[string]any
type Float32List []float32
type UInt32List []uint32
type StringList []string

func (m *JSONObject) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		*m = make(map[string]any)
		return nil
	}

	rc, _, err := argio.OpenReader(value)
	if err != nil {
		return err
	}
	defer rc.Close()

	var tmp map[string]any
	if err := json.NewDecoder(rc).Decode(&tmp); err != nil {
		return pcio.Errorf("failed to parse JSON: %w", err)
	}
	*m = tmp
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
	value = strings.TrimSpace(value)
	if value == "" {
		*m = (*m)[:0]
		return nil
	}

	rc, _, err := argio.OpenReader(value)
	if err != nil {
		return err
	}
	defer rc.Close()

	var arr []float32
	if err := json.NewDecoder(rc).Decode(&arr); err != nil {
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
	value = strings.TrimSpace(value)
	if value == "" {
		*m = (*m)[:0]
		return nil
	}

	rc, _, err := argio.OpenReader(value)
	if err != nil {
		return err
	}
	defer rc.Close()

	var arr []uint32
	if err := json.NewDecoder(rc).Decode(&arr); err != nil {
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
	value = strings.TrimSpace(value)
	if value == "" {
		*m = (*m)[:0]
		return nil
	}

	rc, _, err := argio.OpenReader(value)
	if err != nil {
		return err
	}
	defer rc.Close()

	var arr []string
	if err := json.NewDecoder(rc).Decode(&arr); err != nil {
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
