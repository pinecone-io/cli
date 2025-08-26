package secrets

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
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

var ClientId = configuration.ConfigProperty[string]{
	KeyName:      "client_id",
	ViperStore:   SecretsViper,
	DefaultValue: "",
}

var ClientSecret = configuration.ConfigProperty[string]{
	KeyName:      "client_secret",
	ViperStore:   SecretsViper,
	DefaultValue: "",
}

var ApiKey = configuration.ConfigProperty[string]{
	KeyName:      "api_key",
	ViperStore:   SecretsViper,
	DefaultValue: "",
}

var ProjectAPIKeys = configuration.MarshaledProperty[map[string]string]{
	KeyName:      "project_api_keys",
	ViperStore:   SecretsViper,
	DefaultValue: &map[string]string{},
}

var properties = []configuration.Property{
	ApiKey,
	ClientId,
	ClientSecret,
	OAuth2Token,
	ProjectAPIKeys,
}

var ConfigFile = configuration.ConfigFile{
	FileName:   "secrets",
	FileFormat: "yaml",
	Properties: properties,
	ViperStore: SecretsViper,
}

func init() {
	ConfigFile.Init()

	// Bind environment variables to their associated properties
	SecretsViper.SetEnvPrefix("pinecone")
	_ = SecretsViper.BindEnv(ApiKey.KeyName)
	_ = SecretsViper.BindEnv(ClientId.KeyName)
	_ = SecretsViper.BindEnv(ClientSecret.KeyName)
}

func GetProjectAPIKeys() map[string]string {
	keys := ProjectAPIKeys.Get()
	if keys == nil {
		// if the value is nil, return an empty map to work with
		return map[string]string{}
	}

	return keys
}

func GetProjectAPIKey(projectId string) (string, bool) {
	keys := ProjectAPIKeys.Get()
	if keys == nil {
		return "", false
	}

	key, ok := keys[projectId]
	return key, ok
}

func SetProjectAPIKey(projectId string, apiKey string) {
	keys := ProjectAPIKeys.Get()
	if keys == nil {
		keys = map[string]string{}
	}

	keys[projectId] = apiKey
	ProjectAPIKeys.Set(&keys)
}

func DeleteProjectAPIKey(projectId string) {
	keys := ProjectAPIKeys.Get()
	if keys == nil {
		return
	}
	delete(keys, projectId)
	ProjectAPIKeys.Set(&keys)
}
