package assistants

import (
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_DELETE_ASSISTANT_FILE = "/knowledge/files/%s/%s"
)

type DeleteAssistantFileResponse string

func DeleteKnowledgeFile(kmName string, fileId string) (*DeleteAssistantFileResponse, error) {
	assistantDataUrl, err := GetKnowledgeDataBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.RequestWithoutBodyAndDecode[DeleteAssistantFileResponse](
		assistantDataUrl,
		pcio.Sprintf(URL_DELETE_ASSISTANT_FILE, kmName, fileId),
		http.MethodDelete,
		true,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
