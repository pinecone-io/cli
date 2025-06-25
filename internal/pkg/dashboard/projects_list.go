package dashboard

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_LIST_PROJECTS = "/management/organizations/%s/projects"
)

type Project struct {
	Id                      string `json:"id"`
	Name                    string `json:"name"`
	OrganizationId          string `json:"organization_id"`
	Quota                   string `json:"quota"`
	IndexQuota              string `json:"index_quota"`
	ForceEncryptionWithCmek bool   `json:"force_encryption_with_cmek"`
}

type ProjectsResponse struct {
	Projects []Project `json:"data"`
}

func ListProjects(orgId string) (*ProjectsResponse, error) {
	pineconeApiUrl, err := GetPineconeBaseURL()
	if err != nil {
		return nil, err
	}

	resp, err := network.GetAndDecode[ProjectsResponse](pineconeApiUrl, fmt.Sprintf(URL_LIST_PROJECTS, orgId))
	if err != nil {
		return nil, err
	}

	return resp, nil
}
