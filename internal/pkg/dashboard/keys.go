package dashboard

import (
	"fmt"
)

type KeyResponse struct {
	Keys []Key `json:"keys"`
}

type Key struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	UserLabel     string `json:"user_label"`
	Value         string `json:"value"`
	IntegrationId string `json:"integration_id"`
}

const (
	URL_GET_API_KEYS = "/v2/dashboard/projects/%s/api-keys"
)

func GetApiKeys(project Project) (*KeyResponse, error) {
	url := fmt.Sprintf(URL_GET_API_KEYS, project.GlobalProject.Id)
	return FetchAndDecode[KeyResponse](url, "GET")
}
