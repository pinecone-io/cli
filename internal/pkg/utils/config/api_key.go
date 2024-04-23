package config

import (
	"github.com/spf13/viper"
)

const apiKeyKey = "api_key"

func init() {
	viper.SetDefault(apiKeyKey, "")
}

type ApiKeyConfig struct{}

var ApiKey = ApiKeyConfig{}

func (a ApiKeyConfig) Set(newApiKey string) {
	viper.Set(apiKeyKey, newApiKey)
}

func (a ApiKeyConfig) Get() string {
	return viper.GetString(apiKeyKey)
}
