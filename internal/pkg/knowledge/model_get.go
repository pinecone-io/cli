package knowledge

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_DESCRIBE_KNOWLEDGE_MODEL = "/knowledge/models/%s"
)

func DescribeKnowledgeModel(kmName string) (*KnowledgeModel, error) {
	resp, err := network.GetAndDecode[KnowledgeModel](
		KnowledgeControlPlaneBaseStagingUrl,
		fmt.Sprintf(URL_DESCRIBE_KNOWLEDGE_MODEL, kmName),
	)
	if err != nil {
		return nil, err
	}
	log.Trace().
		Str("model", resp.Name).
		Str("status", string(resp.Status)).
		Str("created_at", resp.CreatedAt).
		Str("updated_at", resp.UpdatedAt).
		Str("metadata", resp.Metadata.ToString()).
		Msg("found model")
	return resp, nil
}
