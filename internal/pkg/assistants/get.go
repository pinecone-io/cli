package assistants

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_DESCRIBE_ASSISTANT = "/assistant/assistants/%s"
)

func DescribeAssistant(name string) (*AssistantModel, error) {
	assistantControlUrl, err := GetAssistantControlBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.GetAndDecode[AssistantModel](
		assistantControlUrl,
		fmt.Sprintf(URL_DESCRIBE_ASSISTANT, name),
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
