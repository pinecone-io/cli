package knowledge

import (
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_DELETE_KNOWLEDGE_FILE = "/knowledge/files/%s/%s"
)

type DeleteKnowledgeFileResponse string

func DeleteKnowledgeFile(kmName string, fileId string) (*DeleteKnowledgeFileResponse, error) {
	resp, err := network.RequestWithoutBodyAndDecode[DeleteKnowledgeFileResponse](
		KnowledgeDataPlaneBaseStagingUrl,
		pcio.Sprintf(URL_DELETE_KNOWLEDGE_FILE, kmName, fileId), http.MethodDelete)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
