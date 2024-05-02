package dashboard

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

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
