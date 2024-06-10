package knowledge

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/environment"
)

func GetKnowledgeDataBaseUrl() (string, error) {
	connectionConfigs, err := environment.GetEnvConfig(config.Environment.Get())
	if err != nil {
		return "", err
	}
	return connectionConfigs.KnowledgeDataPlaneUrl, nil
}

func GetKnowledgeControlBaseUrl() (string, error) {
	connectionConfigs, err := environment.GetEnvConfig(config.Environment.Get())
	if err != nil {
		return "", err
	}
	return connectionConfigs.KnowledgeControlPlaneUrl, nil
}
