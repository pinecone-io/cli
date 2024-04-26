package secrets

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/spf13/viper"
)

var SecretsViper *viper.Viper

const accessTokenKey string = "access_token"
const apiKeyKey string = "api_key"

func init() {
	SecretsViper = viper.New()
	locations := configuration.NewConfigLocations()

	SecretsViper.SetConfigName("secrets") // name of config file (without extension)
	SecretsViper.SetConfigType("yaml")
	SecretsViper.AddConfigPath(locations.ConfigPath)

	SecretsViper.SetDefault(apiKeyKey, "")
	SecretsViper.SetDefault(accessTokenKey, "")
	SecretsViper.SafeWriteConfig()
}

func LoadSecrets() {
	err := SecretsViper.ReadInConfig() // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		exit.Error(err)
	}
}

func SaveSecrets() {
	err := SecretsViper.WriteConfig()
	if err != nil {
		exit.Error(err)
	}
}
