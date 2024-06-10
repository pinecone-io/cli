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

	knowledgeControlUrl, err := GetKnowledgeControlBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.PostAndDecode[CreateKnowledgeModelRequest, KnowledgeModel](
		knowledgeControlUrl,
		URL_POST_KNOWLEDGE_MODEL,
		true,
		body,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
