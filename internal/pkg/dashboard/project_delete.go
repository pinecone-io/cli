package dashboard

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_DELETE_PROJECT = "/v2/dashboard/projects/%s"
)

type DeletePostResponse struct {
	Success bool `json:"success"`
}

func DeleteProject(orgId string, projId string) (*DeletePostResponse, error) {
	path := pcio.Sprintf(URL_DELETE_PROJECT, projId)
	resp, err := DeleteAndDecode[DeletePostResponse](path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
