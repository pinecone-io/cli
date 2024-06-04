package dashboard

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
)

const (
	URL_LIST_ORGANIZATIONS = "/v2/dashboard/organizations"
)

type OrganizationsResponse struct {
	Organizations []Organization `json:"organizations"`
}

type Organization struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Projects []Project `json:"projects"`
}

type Project struct {
	Id            string        `json:"id"`
	Name          string        `json:"name"`
	GlobalProject GlobalProject `json:"globalProject"`
}

type GlobalProject struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Quota      string `json:"quota"`
	IndexQuota string `json:"indexQuota"`
}

func ListOrganizations() (*OrganizationsResponse, error) {
	resp, err := network.GetAndDecode[OrganizationsResponse](DashboardBaseURL, URL_LIST_ORGANIZATIONS)
	if err != nil {
		return nil, err
	}
	for _, org := range resp.Organizations {
		log.Trace().
			Str("org", string(org.Name)).
			Msg("found org")
		for _, proj := range org.Projects {
			log.Trace().
				Str("org", string(org.Name)).
				Str("project", proj.Name).
				Msg("found project in org")
		}
	}
	return resp, nil
}
