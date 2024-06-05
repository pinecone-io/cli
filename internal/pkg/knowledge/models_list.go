package knowledge

import (
	"encoding/json"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_LIST_KNOWLEDGE_MODELS = "/knowledge/models"
)

type KnowledgeModel struct {
	Name      string                    `json:"name"`
	Metadata  KnowledgeMetadata         `json:"metadata"`
	Status    KnowledgeModelStatusState `json:"status"`
	CreatedAt string                    `json:"created_at"`
	UpdatedAt string                    `json:"updated_at"`
}

type KnowledgeMetadata map[string]interface{}

func (kmm *KnowledgeMetadata) ToString() string {
	jsonData, err := json.Marshal(kmm)
	// TODO : handle swallowing decoding error
	if err != nil {
		return ""
	}
	return string(jsonData)
}

type KnowledgeModelStatusState string

const (
	Initializing         KnowledgeModelStatusState = "Initializing"
	InitializationFailed KnowledgeModelStatusState = "Failed"
	Ready                KnowledgeModelStatusState = "Ready"
	ScalingDown          KnowledgeModelStatusState = "Terminating"
)

type ListKnowledgeModelsResponse struct {
	KnowledgeModels []KnowledgeModel `json:"knowledge_models"`
}

func ListKnowledgeModels() (*ListKnowledgeModelsResponse, error) {
	resp, err := network.GetAndDecode[ListKnowledgeModelsResponse](
		GetKnowledgeControlBaseUrl(),
		URL_LIST_KNOWLEDGE_MODELS)
	if err != nil {
		return nil, err
	}
	for _, model := range resp.KnowledgeModels {
		log.Trace().
			Str("model", model.Name).
			Str("status", string(model.Status)).
			Str("created_at", model.CreatedAt).
			Str("updated_at", model.UpdatedAt).
			Str("metadata", model.Metadata.ToString()).
			Msg("found model")
	}
	return resp, nil
}
