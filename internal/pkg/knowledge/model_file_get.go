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
	checkStr := fmt.Sprintf(URL_GET_KNOWLEDGE_MODEL_FILE, kmName, fileId)
	fmt.Printf("CHECK URL STR: %s\n", checkStr)

	resp, err := network.GetAndDecode[KnowledgeFileModel](
		KnowledgeDataPlaneBaseStagingUrl,
		fmt.Sprintf(URL_GET_KNOWLEDGE_MODEL_FILE, kmName, fileId),
	)
	if err != nil {
		exit.Error(err)
	}
	log.Trace().
		Str("name", resp.Name).
		Str("id", resp.Id).
		Str("metadata", resp.Metadata.ToString()).
		Str("mime_type", resp.MimeType).
		Str("created_on", resp.CreatedOn).
		Str("updated_on", resp.UpdatedOn).
		Str("status", string(resp.Status)).
		Msg("found file")

	return resp, nil
}
