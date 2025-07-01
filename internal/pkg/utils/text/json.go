package text

import (
	"encoding/json"
)

func InlineJSON(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func IndentJSON(data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return ""
	}
	return string(jsonData)
}
