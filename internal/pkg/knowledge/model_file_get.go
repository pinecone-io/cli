package knowledge

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_GET_KNOWLEDGE_MODEL_FILE = "/knowledge/files/%s/%s"
)

func DescribeKnowledgeModelFile(kmName string, fileId string) (*KnowledgeFileModel, error) {
	resp, err := network.GetAndDecode[KnowledgeFileModel](
		GetKnowledgeDataBaseUrl(),
		fmt.Sprintf(URL_GET_KNOWLEDGE_MODEL_FILE, kmName, fileId),
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
		Msg("found file")

	return resp, nil
}
