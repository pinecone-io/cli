package knowledge

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_POST_KNOWLEDGE_MODEL = "/knowledge/models"
)

type CreateKnowledgeModelRequest struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

func CreateKnowledgeModel(name string) (*KnowledgeModel, error) {
	body := CreateKnowledgeModelRequest{
		Name: name,
	}
	resp, err := network.PostAndDecode[CreateKnowledgeModelRequest, KnowledgeModel](
		GetKnowledgeControlBaseUrl(),
		URL_POST_KNOWLEDGE_MODEL, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
