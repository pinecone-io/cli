package text

import (
	"encoding/json"
)

func InlineJSON(data any) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func IndentJSON(data any) string {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return ""
	}
	return string(jsonData)
}
