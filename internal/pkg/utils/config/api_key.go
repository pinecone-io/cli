package config

import (
	"github.com/spf13/viper"
)

const apiKeyKey string = "api_key"

func init() {
	viper.SetDefault(apiKeyKey, "")
}

type ApiKeyConfig struct {
	keyName string
}

var ApiKey = ApiKeyConfig{
	keyName: apiKeyKey,
}

func (a ApiKeyConfig) Set(newApiKey string) {
	viper.Set(apiKeyKey, newApiKey)
	SaveConfig()
}

func (a ApiKeyConfig) Get() string {
	return viper.GetString(apiKeyKey)
}
