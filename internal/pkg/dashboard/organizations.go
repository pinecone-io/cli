package dashboard

const (
	URL_GET_ORGANIZATIONS = "/v2/dashboard/organizations"
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

func GetOrganizations() (*OrganizationsResponse, error) {
	return FetchAndDecode[OrganizationsResponse](URL_GET_ORGANIZATIONS, "GET")
}
