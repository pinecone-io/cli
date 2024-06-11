package assistants

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_DESCRIBE_ASSISTANT = "/knowledge/models/%s"
)

func DescribeAssistant(name string) (*AssistantModel, error) {
	knowledgeControlUrl, err := GetKnowledgeControlBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.GetAndDecode[AssistantModel](
		knowledgeControlUrl,
		fmt.Sprintf(URL_DESCRIBE_ASSISTANT, name),
		true,
	)
	if err != nil {
		return nil, err
	}
	log.Trace().
		Str("name", resp.Name).
		Str("status", string(resp.Status)).
		Str("created_at", resp.CreatedAt).
		Str("updated_at", resp.UpdatedAt).
		Str("metadata", resp.Metadata.ToString()).
		Msg("found assistant")
	return resp, nil
}
