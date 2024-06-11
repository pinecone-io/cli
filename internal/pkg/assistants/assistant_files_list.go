package assistants

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_LIST_ASSISTANT_FILES         = "/knowledge/files/%s"
	URL_LIST_ASSISTANT_FILES_STAGING = "/assistant/files/%s"
)

func getListAssistantFilesUrl() string {
	if config.Environment.Get() == "production" {
		return URL_LIST_ASSISTANT_FILES
	} else {
		return URL_LIST_ASSISTANT_FILES_STAGING
	}
}

type AssistantFileModel struct {
	Name      string                   `json:"name"`
	Id        string                   `json:"id"`
	Metadata  AssistantMetadata        `json:"metadata"`
	CreatedOn string                   `json:"created_on"`
	UpdatedOn string                   `json:"updated_on"`
	Status    AssistantFileStatusState `json:"status"`
}

type AssistantFileStatusState string

const (
	Processing AssistantFileStatusState = "Processing"
	Available  AssistantFileStatusState = "Available"
	Deleting   AssistantFileStatusState = "Deleting"
)

type ListAssistantFilesResponse struct {
	Files []AssistantFileModel `json:"files"`
}

func ListAssistantFiles(name string) (*ListAssistantFilesResponse, error) {
	assistantDataUrl, err := GetAssistantDataBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := network.GetAndDecode[ListAssistantFilesResponse](
		assistantDataUrl,
		fmt.Sprintf(getListAssistantFilesUrl(), name),
		true,
	)
	if err != nil {
		return nil, err
	}
	for _, model := range resp.Files {
		log.Trace().
			Str("name", model.Name).
			Str("status", string(model.Status)).
			Str("created_on", model.CreatedOn).
			Str("updated_on", model.UpdatedOn).
			Str("metadata", model.Metadata.ToString()).
			Msg("found assistant")
	}
	return resp, nil
}
