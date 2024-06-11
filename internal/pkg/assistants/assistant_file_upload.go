package assistants

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_ASSISTANT_FILE_UPLOAD = "/knowledge/files/%s"
)

func UploadAssistantFile(name string, filePath string) (*AssistantFileModel, error) {
	assistantDataUrl, err := GetKnowledgeDataBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.PostAndDecodeMultipartFormData[AssistantFileModel](
		assistantDataUrl,
		fmt.Sprintf(URL_ASSISTANT_FILE_UPLOAD, name),
		true,
		filePath,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
