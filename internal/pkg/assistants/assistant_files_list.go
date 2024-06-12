package assistants

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_LIST_ASSISTANT_FILES = "/assistant/files/%s"
)

type AssistantFileModel struct {
	Name      string                   `json:"name"`
	Id        string                   `json:"id"`
	Metadata  AssistantMetadata        `json:"metadata"`
	CreatedOn string                   `json:"created_on"`
	UpdatedOn string                   `json:"updated_on"`
	Status    AssistantFileStatusState `json:"status"`
	Size      int64                    `json:"size"`
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
		fmt.Sprintf(URL_LIST_ASSISTANT_FILES, name),
		true,
	)
	if err != nil {
		return nil, err
	}
	for _, file := range resp.Files {
		log.Trace().
			Str("name", file.Name).
			Str("status", string(file.Status)).
			Str("created_on", file.CreatedOn).
			Str("updated_on", file.UpdatedOn).
			Str("metadata", file.Metadata.ToString()).
			Int64("size", file.Size).
			Msg("found file")
	}
	return resp, nil
}
