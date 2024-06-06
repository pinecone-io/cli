package knowledge

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_KNOWLEDGE_FILE_UPLOAD = "/knowledge/files/%s"
)

func UploadKnowledgeFile(kmName string, filePath string) (*KnowledgeFileModel, error) {
	resp, err := network.PostAndDecodeMultipartFormData[KnowledgeFileModel](
		GetKnowledgeDataBaseUrl(),
		fmt.Sprintf(URL_KNOWLEDGE_FILE_UPLOAD, kmName),
		true,
		filePath,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
