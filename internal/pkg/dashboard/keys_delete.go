package dashboard

import (
	"net/http"

	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_DELETE_API_KEYS = "/v2/dashboard/projects/%s/api-key"
)

type DeleteApiKeyRequest struct {
	Label    string `json:"label"`
	UserName string `json:"userName"`
}

type DeleteApiKeyResponse struct {
	Success bool `json:"success"`
}

func DeleteApiKey(projId string, key Key) (*DeleteApiKeyResponse, error) {
	path := pcio.Sprintf(URL_DELETE_API_KEYS, projId)
	body := DeleteApiKeyRequest{
		Label:    key.UserLabel,
		UserName: key.UserName,
	}

	dashboardUrl, err := GetDashboardBaseURL()
	if err != nil {
		return nil, err
	}

	resp, err := network.RequestWithBodyAndDecode[DeleteApiKeyRequest, DeleteApiKeyResponse](
		dashboardUrl,
		path,
		http.MethodDelete,
		false,
		body,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
