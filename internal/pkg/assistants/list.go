package assistants

import (
	"encoding/json"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_LIST_ASSISTANTS = "/assistant/assistants"
)

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
	if err != nil {
		return "ERROR: could not parse AssistantMetadata to string"
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
	Assistants []AssistantModel `json:"assistants"`
}

func ListAssistants() (*ListAssistantsResponse, error) {
	assistantControlUrl, err := GetAssistantControlBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.GetAndDecode[ListAssistantsResponse](
		assistantControlUrl,
		URL_LIST_ASSISTANTS,
		true,
	)
	if err != nil {
		return nil, err
	}
	for _, model := range resp.Assistants {
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
