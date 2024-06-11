package assistants

import (
	"encoding/json"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_LIST_ASSISTANTS         = "/knowledge/models"
	URL_LIST_ASSISTANTS_STAGING = "/assistant/assistants"
)

func getListAssistantsUrl() string {
	if config.Environment.Get() == "production" {
		return URL_LIST_ASSISTANTS
	} else {
		return URL_LIST_ASSISTANTS_STAGING
	}
}

type AssistantModel struct {
	Name      string               `json:"name"`
	Metadata  AssistantMetadata    `json:"metadata"`
	Status    AssistantStatusState `json:"status"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
}

type AssistantMetadata map[string]interface{}

func (kmm *AssistantMetadata) ToString() string {
	jsonData, err := json.Marshal(kmm)
	// TODO : handle swallowing decoding error
	if err != nil {
		return ""
	}
	return string(jsonData)
}

type AssistantStatusState string

const (
	Initializing AssistantStatusState = "Initializing"
	Failed       AssistantStatusState = "Failed"
	Ready        AssistantStatusState = "Ready"
	Terminating  AssistantStatusState = "Terminating"
)

type ListAssistantsResponse struct {
	KnowledgeModels []AssistantModel `json:"knowledge_models"`
}

func ListAssistants() (*ListAssistantsResponse, error) {
	assistantControlUrl, err := GetAssistantControlBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.GetAndDecode[ListAssistantsResponse](
		assistantControlUrl,
		getListAssistantsUrl(),
		true,
	)
	if err != nil {
		return nil, err
	}
	for _, model := range resp.KnowledgeModels {
		log.Trace().
			Str("name", model.Name).
			Str("status", string(model.Status)).
			Str("created_at", model.CreatedAt).
			Str("updated_at", model.UpdatedAt).
			Str("metadata", model.Metadata.ToString()).
			Msg("found assistant")
	}
	return resp, nil
}
