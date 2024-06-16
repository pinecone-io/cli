package assistants

import (
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_DELETE_ASSISTANT_FILE = "/assistant/files/%s/%s"
)

type DeleteAssistantFileResponse string

func DeleteAssistantFile(asstName string, fileId string) (*DeleteAssistantFileResponse, error) {
	assistantDataUrl, err := GetAssistantDataBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.RequestWithoutBodyAndDecode[DeleteAssistantFileResponse](
		assistantDataUrl,
		pcio.Sprintf(URL_DELETE_ASSISTANT_FILE, asstName, fileId),
		http.MethodDelete,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
