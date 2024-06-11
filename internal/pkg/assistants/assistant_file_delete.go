package assistants

import (
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_DELETE_ASSISTANT_FILE         = "/knowledge/files/%s/%s"
	URL_DELETE_ASSISTANT_FILE_STAGING = "/assistant/files/%s/%s"
)

func getDeleteAssistantFileUrl() string {
	if config.Environment.Get() == "production" {
		return URL_DELETE_ASSISTANT_FILE
	} else {
		return URL_DELETE_ASSISTANT_FILE_STAGING
	}
}

type DeleteAssistantFileResponse string

func DeleteKnowledgeFile(kmName string, fileId string) (*DeleteAssistantFileResponse, error) {
	assistantDataUrl, err := GetAssistantDataBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.RequestWithoutBodyAndDecode[DeleteAssistantFileResponse](
		assistantDataUrl,
		pcio.Sprintf(getDeleteAssistantFileUrl(), kmName, fileId),
		http.MethodDelete,
		true,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
