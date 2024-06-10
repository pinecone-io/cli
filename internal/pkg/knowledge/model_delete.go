package knowledge

import (
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_DELETE_KNOWLEDGE_MODEL = "/knowledge/models/%s"
)

type DeleteKnowledgeModelResponse struct {
	Success bool `json:"success"`
}

func DeleteKnowledgeModel(kmName string) (*DeleteKnowledgeModelResponse, error) {

	knowledgeControlUrl, err := GetKnowledgeControlBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.RequestWithoutBodyAndDecode[DeleteKnowledgeModelResponse](
		knowledgeControlUrl,
		pcio.Sprintf(URL_DELETE_KNOWLEDGE_MODEL, kmName),
		http.MethodDelete,
		true,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
