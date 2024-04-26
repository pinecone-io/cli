package secrets

type AccessTokenConfig struct {
	keyName string
}

var AccessToken = AccessTokenConfig{
	keyName: accessTokenKey,
}

func (a AccessTokenConfig) Set(newApiKey string) {
	SecretsViper.Set(accessTokenKey, newApiKey)
	SaveSecrets()
}

func (a AccessTokenConfig) Get() string {
	return SecretsViper.GetString(accessTokenKey)
}
