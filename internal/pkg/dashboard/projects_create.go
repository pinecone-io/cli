package dashboard

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_POST_PROJECTS = "/v2/dashboard/organizations/%s/projects"
)

type CreateProjectRequest struct {
	Name        string `json:"name"`
	PodQuota    int32  `json:"quota"`
	Environment string `json:"environment"`
}

type CreateProjectResponse struct {
	Success bool    `json:"success"`
	Project Project `json:"globalProject"`
}

func CreateProject(orgId string, projName string, podQuota int32) (*CreateProjectResponse, error) {
	path := pcio.Sprintf(URL_POST_PROJECTS, orgId)
	body := CreateProjectRequest{
		Name:        projName,
		PodQuota:    podQuota,
		Environment: "serverless",
	}

	dashboardUrl, err := GetDashboardBaseURL()
	if err != nil {
		return nil, err
	}

	resp, err := network.PostAndDecode[CreateProjectRequest, CreateProjectResponse](
		dashboardUrl,
		path,
		body,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
