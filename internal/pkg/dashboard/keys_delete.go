package dashboard

import (
	"net/http"

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

	resp, err := RequestWithBodyAndDecode[DeleteApiKeyRequest, DeleteApiKeyResponse](path, http.MethodDelete, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
