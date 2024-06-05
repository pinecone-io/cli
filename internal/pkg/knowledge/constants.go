package knowledge

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
)

const KnowledgeDataPlaneBaseStagingUrl = "https://staging-data.ke.pinecone.io"
const KnowledgeDataPlaneBaseUrl = "https://prod-1-data.ke.pinecone.io"

const KnowledgeControlPlaneBaseStagingUrl = "https://api-staging.pinecone.io"
const KnowledgeControlPlaneBaseUrl = "https://api.pinecone.io"

func GetKnowledgeDataBaseUrl() string {
	if config.Staging.Get() {
		return KnowledgeDataPlaneBaseStagingUrl
	} else {
		return KnowledgeDataPlaneBaseUrl
	}
}

func GetKnowledgeControlBaseUrl() string {
	if config.Staging.Get() {
		return KnowledgeControlPlaneBaseStagingUrl
	} else {
		return KnowledgeControlPlaneBaseUrl
	}
}
