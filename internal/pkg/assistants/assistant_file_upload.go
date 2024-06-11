package assistants

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_ASSISTANT_FILE_UPLOAD         = "/knowledge/files/%s"
	URL_ASSISTANT_FILE_UPLOAD_STAGING = "/assistant/files/%s"
)

func getAssistantFileUploadUrl() string {
	if config.Environment.Get() == "production" {
		return URL_ASSISTANT_FILE_UPLOAD
	} else {
		return URL_ASSISTANT_FILE_UPLOAD_STAGING
	}
}

func UploadAssistantFile(name string, filePath string) (*AssistantFileModel, error) {
	assistantDataUrl, err := GetAssistantDataBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.PostAndDecodeMultipartFormData[AssistantFileModel](
		assistantDataUrl,
		fmt.Sprintf(getAssistantFileUploadUrl(), name),
		true,
		filePath,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
