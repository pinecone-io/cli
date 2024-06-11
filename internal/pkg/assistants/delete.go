package assistants

import (
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_DELETE_ASSISTANT         = "/knowledge/models/%s"
	URL_DELETE_ASSISTANT_STAGING = "/assistant/assistants/%s"
)

func getDeleteAssistantUrl() string {
	if config.Environment.Get() == "production" {
		return URL_DELETE_ASSISTANT
	} else {
		return URL_DELETE_ASSISTANT_STAGING
	}
}

type DeleteKnowledgeModelResponse struct {
	Success bool `json:"success"`
}

func DeleteAssistant(name string) (*DeleteKnowledgeModelResponse, error) {

	assistantControlUrl, err := GetAssistantControlBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.RequestWithoutBodyAndDecode[DeleteKnowledgeModelResponse](
		assistantControlUrl,
		pcio.Sprintf(getDeleteAssistantUrl(), name),
		http.MethodDelete,
		true,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
