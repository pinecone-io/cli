package knowledge

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_LIST_KNOWLEDGE_MODEL_FILES = "/knowledge/files/%s"
)

type KnowledgeFileModel struct {
	Name      string                   `json:"name"`
	Id        string                   `json:"id"`
	Metadata  KnowledgeMetadata        `json:"metadata"`
	CreatedOn string                   `json:"created_on"`
	UpdatedOn string                   `json:"updated_on"`
	Status    KnowledgeFileStatusState `json:"status"`
}

type KnowledgeFileStatusState string

const (
	Processing KnowledgeFileStatusState = "Processing"
	Available  KnowledgeFileStatusState = "Available"
	Deleting   KnowledgeFileStatusState = "Deleting"
)

type ListKnowledgeModelFilesResponse struct {
	Files []KnowledgeFileModel `json:"files"`
}

func ListKnowledgeModelFiles(kmName string) (*ListKnowledgeModelFilesResponse, error) {
	resp, err := network.GetAndDecode[ListKnowledgeModelFilesResponse](
		GetKnowledgeDataBaseUrl(),
		fmt.Sprintf(URL_LIST_KNOWLEDGE_MODEL_FILES, kmName),
		true,
	)
	if err != nil {
		return nil, err
	}
	for _, model := range resp.Files {
		log.Trace().
			Str("model", model.Name).
			Str("status", string(model.Status)).
			Str("created_on", model.CreatedOn).
			Str("updated_on", model.UpdatedOn).
			Str("metadata", model.Metadata.ToString()).
			Msg("found model")
	}
	return resp, nil
}
