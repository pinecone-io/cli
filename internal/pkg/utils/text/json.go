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

func PrettyPrintJSON(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	pcio.Println(string(jsonData))
	return string(jsonData), nil
}
