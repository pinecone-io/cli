package text

import (
	"encoding/json"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
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
	pcio.Println(string(jsonData))
	return string(jsonData)
}
