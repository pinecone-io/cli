package dashboard

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

type KeyResponse struct {
	Keys []Key `json:"keys"`
}

type Key struct {
	Id        string `json:"id"`
	UserLabel string `json:"user_label"`
	Value     string `json:"value"`
	UserName  string `json:"user_name"`
}

const (
	URL_GET_API_KEYS = "/v2/dashboard/projects/%s/api-keys"
)

func GetApiKeys(project GlobalProject) (*KeyResponse, error) {
	return GetApiKeysById(project.Id)
}

func GetApiKeysById(projectId string) (*KeyResponse, error) {
	dashboardUrl, err := GetDashboardBaseURL()
	if err != nil {
		return nil, err
	}

	url := pcio.Sprintf(URL_GET_API_KEYS, projectId)
	return network.GetAndDecode[KeyResponse](dashboardUrl, url)
}
