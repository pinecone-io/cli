package assistants

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_DESCRIBE_ASSISTANT_FILE         = "/knowledge/files/%s/%s"
	URL_DESCRIBE_ASSISTANT_FILE_STAGING = "/assistant/files/%s/%s"
)

func getDescribeAssistantFileUrl() string {
	if config.Environment.Get() == "production" {
		return URL_DESCRIBE_ASSISTANT_FILE
	} else {
		return URL_DESCRIBE_ASSISTANT_FILE_STAGING
	}
}

func DescribeAssistantFile(name string, fileId string) (*AssistantFileModel, error) {
	assistantDataUrl, err := GetAssistantDataBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.GetAndDecode[AssistantFileModel](
		assistantDataUrl,
		fmt.Sprintf(getDescribeAssistantFileUrl(), name, fileId),
		true,
	)
	if err != nil {
		exit.Error(err)
	}
	log.Trace().
		Str("name", resp.Name).
		Str("id", resp.Id).
		Str("metadata", resp.Metadata.ToString()).
		Str("created_on", resp.CreatedOn).
		Str("updated_on", resp.UpdatedOn).
		Str("status", string(resp.Status)).
		Int64("size", resp.Size).
		Msg("found file")

	return resp, nil
}
