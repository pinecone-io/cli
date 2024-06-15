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

func (am *AssistantMetadata) ToString() string {
	jsonData, err := json.Marshal(am)
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
		false,
	)
	if err != nil {
		return nil, err
	}
	for _, assistant := range resp.Assistants {
		log.Trace().
			Str("name", assistant.Name).
			Str("status", string(assistant.Status)).
			Str("created_at", assistant.CreatedAt).
			Str("updated_at", assistant.UpdatedAt).
			Str("metadata", assistant.Metadata.ToString()).
			Msg("found assistant")
	}
	return resp, nil
}
