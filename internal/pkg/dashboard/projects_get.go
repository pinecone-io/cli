package dashboard

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

func GetProjectByName(orgName string, projName string) (*Project, error) {
	orgs, err := ListOrganizations()
	if err != nil {
		return nil, err
	}
	for _, org := range orgs.Organizations {
		if org.Name == orgName {
			for _, proj := range org.Projects {
				if proj.Name == projName {
					return &proj, nil
				}
			}
		}
	}
	return nil, error(pcio.Errorf("project name %s not found in organization %s", projName, orgName))
}

func GetProjectById(orgId string, projId string) (*Project, error) {
	orgs, err := ListOrganizations()
	if err != nil {
		return nil, err
	}

	for _, org := range orgs.Organizations {
		if org.Id == orgId {
			for _, proj := range org.Projects {
				if proj.Id == projId {
					return &proj, nil
				}
			}
		}
	}
	return nil, error(pcio.Errorf("project id %s not found in organization %s", projId, orgId))
}
