package dashboard

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	pc_oauth2 "github.com/pinecone-io/cli/internal/pkg/utils/oauth2"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/rs/zerolog/log"
)

const (
	URL_LIST_ORGANIZATIONS = "/v2/dashboard/organizations"
)

type OrganizationsResponse struct {
	Organizations []Organization `json:"newOrgs"`
}

type Organization struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	Projects *[]Project `json:"projects"`
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

	accessToken := secrets.OAuth2Token.Get()
	claims, err := pc_oauth2.ParseClaimsUnverified(&accessToken)
	if err != nil {
		return nil, err
	}

	// Match the organization to the jwt token's orgId if possible
	for i := range resp.Organizations {
		if resp.Organizations[i].Id == claims.OrgId {
			org := &resp.Organizations[i]
			projects, err := ListProjects(org.Id)
			if err != nil {
				log.Err(err).
					Msg(fmt.Sprintf("Error listing projects for organization %s. Please create an organization if your account is not associated with one.", org.Name))
				pcio.Printf("Error listing projects for organization %s: %s\n", org.Name, err)
			} else {
				org.Projects = &projects.Projects
			}
		} else {
			log.Error()
		}
	}

	return resp, nil
}
