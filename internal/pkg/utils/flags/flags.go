package flags

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

type JSONObject map[string]any

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
	for k, v := range tmp {
		(*m)[k] = v
	}
	return nil
}

func (m *JSONObject) String() string {
	if m == nil || len(*m) == 0 {
		return ""
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func (*JSONObject) Type() string { return "json" }
