package dashboard

import (
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
	Success       bool          `json:"success"`
	GlobalProject GlobalProject `json:"globalProject"`
}

func GetProjectByName(orgName string, projName string) (*GlobalProject, error) {
	orgs, err := ListOrganizations()
	if err != nil {
		return nil, err
	}
	for _, org := range orgs.Organizations {
		if org.Name == orgName {
			for _, proj := range org.Projects {
				if proj.GlobalProject.Name == projName {
					return &proj.GlobalProject, nil
				}
			}
		}
	}
	return nil, error(pcio.Errorf("project name %s not found in organization %s", projName, orgName))
}

func GetProjectById(orgId string, projId string) (*GlobalProject, error) {
	orgs, err := ListOrganizations()
	if err != nil {
		return nil, err
	}
	for _, org := range orgs.Organizations {
		if org.Id == orgId {
			for _, proj := range org.Projects {
				if proj.GlobalProject.Id == projId {
					return &proj.GlobalProject, nil
				}
			}
		}
	}
	return nil, error(pcio.Errorf("project id %s not found in organization %s", projId, orgId))
}

func CreateProject(orgId string, projName string, podQuota int32) (*CreateProjectResponse, error) {
	path := pcio.Sprintf(URL_POST_PROJECTS, orgId)
	body := CreateProjectRequest{
		Name:        projName,
		PodQuota:    podQuota,
		Environment: "serverless",
	}
	resp, err := PostAndDecode[CreateProjectRequest, CreateProjectResponse](path, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
