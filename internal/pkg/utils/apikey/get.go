package apikey

import (
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
)

func GetApiKey() string {
	storedApiKey := secrets.ApiKey.Get()
	envApiKey := os.Getenv("PINECONE_API_KEY")

	if storedApiKey != "" {
		return storedApiKey
	}

	return envApiKey
}
