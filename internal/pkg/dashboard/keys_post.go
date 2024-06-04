package dashboard

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/network"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

const (
	URL_POST_API_KEYS = "/v2/dashboard/projects/%s/api-key"
)

type CreateApiKeyRequest struct {
	Label string `json:"label"`
}

type CreateApiKeyResponse struct {
	Success bool `json:"success"`
	Key     Key  `json:"key"`
}

func CreateApiKey(projId string, keyName string) (*CreateApiKeyResponse, error) {
	path := pcio.Sprintf(URL_POST_API_KEYS, projId)
	body := CreateApiKeyRequest{
		Label: keyName,
	}
	resp, err := network.PostAndDecode[CreateApiKeyRequest, CreateApiKeyResponse](DashboardBaseURL, path, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
