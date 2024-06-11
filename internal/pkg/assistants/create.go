package assistants

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_CREATE_ASSISTANT = "/knowledge/models"
)

type CreateAssistantRequest struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

func CreateAssistant(name string) (*AssistantModel, error) {
	body := CreateAssistantRequest{
		Name: name,
	}

	knowledgeControlUrl, err := GetKnowledgeControlBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.PostAndDecode[CreateAssistantRequest, AssistantModel](
		knowledgeControlUrl,
		URL_CREATE_ASSISTANT,
		true,
		body,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
