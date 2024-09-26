package secrets

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var SecretsViper *viper.Viper = viper.New()

var OAuth2Token = configuration.MarshaledProperty[oauth2.Token]{
	KeyName:    "oauth2_token",
	ViperStore: SecretsViper,
	DefaultValue: &oauth2.Token{
		AccessToken: "",
	},
}
var (
	ApiKey = configuration.ConfigProperty[string]{
		KeyName:      "api_key",
		ViperStore:   SecretsViper,
		DefaultValue: "",
	}
)
var properties = []configuration.Property{
	ApiKey,
	OAuth2Token,
}

var ConfigFile = configuration.ConfigFile{
	FileName:   "secrets",
	FileFormat: "yaml",
	Properties: properties,
	ViperStore: SecretsViper,
}

func init() {
	ConfigFile.Init()

	SecretsViper.SetEnvPrefix("pinecone")
	err := SecretsViper.BindEnv(ApiKey.KeyName)
	if err != nil {
		exit.Error(err)
	}
}
