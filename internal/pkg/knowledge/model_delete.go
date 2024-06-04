package knowledge

import (
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_DELETE_KNOWLEDGE_MODEL = "/knowledge/models/%s"
)

type DeleteKnowledgeModelRequest struct {
	Name string `json:"name"`
}

type DeleteKnowledgeModelResponse struct {
	Success bool `json:"success"`
}

func DeleteKnowledgeModel(kmName string) (*DeleteKnowledgeModelResponse, error) {

	resp, err := network.RequestWithoutBodyAndDecode[DeleteKnowledgeModelResponse](
		KnowledgeControlPlaneBaseStagingUrl,
		pcio.Sprintf(URL_DELETE_KNOWLEDGE_MODEL, kmName), http.MethodDelete)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
