package secrets

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/spf13/viper"
)

var SecretsViper *viper.Viper

const accessTokenKey string = "access_token"
const refreshTokenKey string = "refresh_token"
const apiKeyKey string = "api_key"

type ConfigProperty struct {
	KeyName string
}

func (c ConfigProperty) Set(value string) {
	SecretsViper.Set(c.KeyName, value)
	SaveSecrets()
}

func (c ConfigProperty) Get() string {
	return SecretsViper.GetString(c.KeyName)
}

var (
	RefreshToken = ConfigProperty{KeyName: refreshTokenKey}
	AccessToken  = ConfigProperty{KeyName: accessTokenKey}
	ApiKey       = ConfigProperty{KeyName: apiKeyKey}
)

func init() {
	SecretsViper = viper.New()
	locations := configuration.NewConfigLocations()

	SecretsViper.SetConfigName("secrets") // name of config file (without extension)
	SecretsViper.SetConfigType("yaml")
	SecretsViper.AddConfigPath(locations.ConfigPath)

	SecretsViper.SetDefault(apiKeyKey, "")
	SecretsViper.SetDefault(accessTokenKey, "")
	SecretsViper.SetDefault(refreshTokenKey, "")
	SecretsViper.SafeWriteConfig()
}

func Clear() {
	ApiKey.Set("")
	AccessToken.Set("")
	RefreshToken.Set("")
	SaveSecrets()
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
