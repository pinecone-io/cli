package assistants

import (
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_DELETE_ASSISTANT = "/assistant/assistants/%s"
)

type DeleteAssistantResponse struct {
	Success bool `json:"success"`
}

func DeleteAssistant(name string) (*DeleteAssistantResponse, error) {

	assistantControlUrl, err := GetAssistantControlBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.RequestWithoutBodyAndDecode[DeleteAssistantResponse](
		assistantControlUrl,
		pcio.Sprintf(URL_DELETE_ASSISTANT, name),
		http.MethodDelete,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
