package text

import (
	"bytes"
	"encoding/json"
	"strings"
)

// encode marshals data to JSON with HTML escaping disabled.
// json.Marshal and json.MarshalIndent escape &, <, > as \uXXXX by default —
// a safety measure for embedding JSON in HTML that is incorrect for CLI output.
func encode(data any, indent bool) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if indent {
		enc.SetIndent("", "    ")
	}
	if err := enc.Encode(data); err != nil {
		return ""
	}
	// json.Encoder.Encode appends a trailing newline; trim it so callers
	// control their own newlines (consistent with the old MarshalIndent behavior).
	return strings.TrimRight(buf.String(), "\n")
}

func InlineJSON(data any) string {
	return encode(data, false)
}

func IndentJSON(data any) string {
	return encode(data, true)
}
