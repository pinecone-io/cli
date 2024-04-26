package secrets

type ApiKeyConfig struct {
	keyName string
}

var ApiKey = ApiKeyConfig{
	keyName: apiKeyKey,
}

func (a ApiKeyConfig) Set(newApiKey string) {
	SecretsViper.Set(apiKeyKey, newApiKey)
	SaveSecrets()
}

func (a ApiKeyConfig) Get() string {
	return SecretsViper.GetString(apiKeyKey)
}
