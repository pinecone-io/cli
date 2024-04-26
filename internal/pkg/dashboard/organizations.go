package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetOrganizations(accessToken string) (*OrganizationsResponse, error) {
	url := "https://console-api.pinecone.io/v2/dashboard/organizations"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
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

	var OrganizationsResponse OrganizationsResponse
	err = json.NewDecoder(resp.Body).Decode(&OrganizationsResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &OrganizationsResponse, nil
}
