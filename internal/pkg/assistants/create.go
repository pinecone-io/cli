package assistants

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_CREATE_ASSISTANT         = "/knowledge/models"
	URL_CREATE_ASSISTANT_STAGING = "/assistant/assistants"
)

func getCreateAssistantUrl() string {
	if config.Environment.Get() == "production" {
		return URL_CREATE_ASSISTANT
	} else {
		return URL_CREATE_ASSISTANT_STAGING
	}
}

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
		getCreateAssistantUrl(),
		true,
		body,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
