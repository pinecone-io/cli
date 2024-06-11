package assistants

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_DESCRIBE_ASSISTANT         = "/knowledge/models/%s"
	URL_DESCRIBE_ASSISTANT_STAGING = "/assistant/assistants/%s"
)

func getDescribeAssistantUrl() string {
	if config.Environment.Get() == "production" {
		return URL_DESCRIBE_ASSISTANT
	} else {
		return URL_DESCRIBE_ASSISTANT_STAGING
	}
}

func DescribeAssistant(name string) (*AssistantModel, error) {
	assistantControlUrl, err := GetAssistantControlBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.GetAndDecode[AssistantModel](
		assistantControlUrl,
		fmt.Sprintf(getDescribeAssistantUrl(), name),
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
