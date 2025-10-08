package secrets

import (
	"sync"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var SecretsViper *viper.Viper = viper.New()

var oauth2TokenMu sync.RWMutex
var oAuth2Token = configuration.MarshaledProperty[oauth2.Token]{
	KeyName:      "oauth2_token",
	ViperStore:   SecretsViper,
	DefaultValue: oauth2.Token{},
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

var DefaultAPIKey = configuration.ConfigProperty[string]{
	KeyName:      "api_key",
	ViperStore:   SecretsViper,
	DefaultValue: "",
}

var managedKeysMu sync.RWMutex
var ManagedAPIKeys = configuration.MarshaledProperty[map[string]ManagedKey]{
	KeyName:      "project_api_keys",
	ViperStore:   SecretsViper,
	DefaultValue: map[string]ManagedKey{},
}

var properties = []configuration.Property{
	DefaultAPIKey,
	ClientId,
	ClientSecret,
	oAuth2Token,
	ManagedAPIKeys,
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
	_ = SecretsViper.BindEnv(DefaultAPIKey.KeyName)
	_ = SecretsViper.BindEnv(ClientId.KeyName)
	_ = SecretsViper.BindEnv(ClientSecret.KeyName)
}

func GetOAuth2Token() oauth2.Token {
	oauth2TokenMu.RLock()
	defer oauth2TokenMu.RUnlock()
	return oAuth2Token.Get()
}

func SetOAuth2Token(token oauth2.Token) {
	oauth2TokenMu.Lock()
	defer oauth2TokenMu.Unlock()
	oAuth2Token.Set(token)
}

func ClearOAuth2Token() {
	oauth2TokenMu.Lock()
	defer oauth2TokenMu.Unlock()
	oAuth2Token.Clear()
}

func GetManagedProjectKeys() map[string]ManagedKey {
	managedKeysMu.RLock()
	defer managedKeysMu.RUnlock()

	keys := ManagedAPIKeys.Get()
	if keys == nil {
		// if the value is nil, return an empty map to work with
		return map[string]ManagedKey{}
	}

	return keys
}

func GetProjectManagedKey(projectId string) (ManagedKey, bool) {
	managedKeysMu.RLock()
	defer managedKeysMu.RUnlock()

	keys := ManagedAPIKeys.Get()
	if keys == nil {
		return ManagedKey{}, false
	}

	key, ok := keys[projectId]
	return key, ok
}

func SetProjectManagedKey(managedKey ManagedKey) {
	managedKeysMu.Lock()
	defer managedKeysMu.Unlock()

	keys := ManagedAPIKeys.Get()
	if keys == nil {
		keys = map[string]ManagedKey{}
	}

	keys[managedKey.ProjectId] = managedKey
	ManagedAPIKeys.Set(keys)
}

func DeleteProjectManagedKey(projectId string) {
	managedKeysMu.Lock()
	defer managedKeysMu.Unlock()

	keys := ManagedAPIKeys.Get()
	if keys == nil {
		return
	}
	delete(keys, projectId)
	ManagedAPIKeys.Set(keys)
}

func ClearManagedProjectKeys() {
	managedKeysMu.Lock()
	defer managedKeysMu.Unlock()
	ManagedAPIKeys.Clear()
}

type ManagedKeyOrigin string

const (
	OriginCLICreated  ManagedKeyOrigin = "cli_created"
	OriginUserCreated ManagedKeyOrigin = "user_created"
)

// ManagedKey represents an API key that is being actively managed by the CLI
// Either the CLI created it to work with a project, or a user created it and stored it explicitly
type ManagedKey struct {
	Name           string           `json:"name,omitempty"`
	Id             string           `json:"id,omitempty"`
	Value          string           `json:"value,omitempty"`
	Origin         ManagedKeyOrigin `json:"origin,omitempty"`
	ProjectId      string           `json:"project_id,omitempty"`
	ProjectName    string           `json:"project_name,omitempty"`
	OrganizationId string           `json:"organization_id,omitempty"`
}
