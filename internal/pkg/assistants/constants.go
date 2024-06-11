package assistants

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/environment"
)

func GetAssistantDataBaseUrl() (string, error) {
	connectionConfigs, err := environment.GetEnvConfig(config.Environment.Get())
	if err != nil {
		return "", err
	}
	return connectionConfigs.AssistantDataPlaneUrl, nil
}

func GetAssistantControlBaseUrl() (string, error) {
	connectionConfigs, err := environment.GetEnvConfig(config.Environment.Get())
	if err != nil {
		return "", err
	}
	return connectionConfigs.AssistantControlPlaneUrl, nil
}
