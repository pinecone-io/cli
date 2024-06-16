package assistants

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_CREATE_ASSISTANT = "/assistant/assistants"
)

type CreateAssistantRequest struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

func CreateAssistant(name string) (*AssistantModel, error) {
	body := CreateAssistantRequest{
		Name: name,
	}

	assistantControlUrl, err := GetAssistantControlBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.PostAndDecode[CreateAssistantRequest, AssistantModel](
		assistantControlUrl,
		URL_CREATE_ASSISTANT,
		body,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
