package dashboard

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/environment"
)

func GetDashboardBaseURL() (string, error) {
	connectionConfigs, err := environment.GetEnvConfig(config.Environment.Get())
	if err != nil {
		return "", err
	}
	return connectionConfigs.DashboardUrl, nil
}

func GetPineconeBaseURL() (string, error) {
	connectionConfigs, err := environment.GetEnvConfig(config.Environment.Get())
	if err != nil {
		return "", err
	}
	return connectionConfigs.IndexControlPlaneUrl, nil
}
