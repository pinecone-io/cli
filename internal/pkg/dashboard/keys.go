package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func GetApiKeys(project Project, accessToken string) (*KeyResponse, error) {
	url := fmt.Sprintf("https://console-api.pinecone.io/v2/dashboard/projects/%s/api-keys", project.GlobalProject.Id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error creating request:", err)
		return nil, err
	}

	// Add headers to the request
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("User-Agent", "Pinecone CLI")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, resp.Status)
	}

	var KeyResponse KeyResponse
	err = json.NewDecoder(resp.Body).Decode(&KeyResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &KeyResponse, nil
}
