package dashboard

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_GET_ORGANIZATION = "/v2/dashboard/organizations/%s"
)

type UserRole struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type DescribeOrganizationResponse struct {
	Success   bool       `json:"success"`
	UserRoles []UserRole `json:"userRoles"`
}

func DescribeOrganization(orgId string) (*DescribeOrganizationResponse, error) {
	path := pcio.Sprintf(URL_GET_ORGANIZATION, orgId)
	dashboardUrl, err := GetDashboardBaseURL()
	if err != nil {
		return nil, err
	}
	resp, err := network.GetAndDecode[DescribeOrganizationResponse](dashboardUrl, path)
	if err != nil {
		return nil, err
	}
	for _, userRole := range resp.UserRoles {
		log.Trace().
			Str("email", userRole.Email).
			Str("role", userRole.Role).
			Msg("found user role")
	}
	return resp, nil
}
