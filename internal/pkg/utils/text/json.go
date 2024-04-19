package text

import (
	"fmt"
	"encoding/json"
)

func PrettyPrintJSON(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}