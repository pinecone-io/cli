package dashboard

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_LIST_ORGANIZATIONS = "/v2/dashboard/organizations"
)

type OrganizationsResponse struct {
	Organizations []Organization `json:"newOrgs"`
	Projects      []Project      `json:"projects"`
}

type Organization struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Projects []Project `json:"projects"`
}

type Project struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	OrganizationId string `json:"organization_id"`
	Quota          string `json:"quota"`
	IndexQuota     string `json:"index_quota"`
}

func ListOrganizations() (*OrganizationsResponse, error) {
	dashboardUrl, err := GetDashboardBaseURL()
	if err != nil {
		return nil, err
	}

	resp, err := network.GetAndDecode[OrganizationsResponse](dashboardUrl, URL_LIST_ORGANIZATIONS)
	if err != nil {
		return nil, err
	}

	orgToProjectsMap := make(map[string][]Project)

	// Organize projects into orgs to match the older data structure
	for _, project := range resp.Projects {
		if projects, ok := orgToProjectsMap[project.OrganizationId]; ok {
			projects = append(projects, project)
			orgToProjectsMap[project.OrganizationId] = projects
		} else {
			orgToProjectsMap[project.OrganizationId] = []Project{project}
		}
	}

	// Modify the response to nest projects under their orgs
	for i := range resp.Organizations {
		org := &resp.Organizations[i]

		log.Trace().
			Str("org", string(org.Name)).
			Msg("found org")

		org.Projects = orgToProjectsMap[org.Id]

		for _, proj := range org.Projects {
			log.Trace().
				Str("org", string(org.Name)).
				Str("project", proj.Name).
				Msg("found project in org")
		}
	}

	return resp, nil
}
