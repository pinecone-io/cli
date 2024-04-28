package secrets

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var SecretsViper *viper.Viper = viper.New()

var OAuth2Token = configuration.MarshaledProperty[oauth2.Token]{
	KeyName:      "oauth2_token",
	ViperStore:   SecretsViper,
	DefaultValue: &oauth2.Token{},
}
var (
	ApiKey = configuration.ConfigProperty[string]{
		KeyName:    "api_key",
		ViperStore: SecretsViper,
		// DefaultValue: "",
	}
)
var properties = []configuration.Property{
	ApiKey,
	OAuth2Token,
}

func init() {
	locations := configuration.NewConfigLocations()

	SecretsViper.SetConfigName("secrets") // name of config file (without extension)
	SecretsViper.SetConfigType("yaml")
	SecretsViper.AddConfigPath(locations.ConfigPath)

	for _, property := range properties {
		property.Init()
	}
	SecretsViper.SafeWriteConfig()
}

func Clear() {
	for _, property := range properties {
		property.Clear()
	}
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
