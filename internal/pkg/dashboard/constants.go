package dashboard

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
)

const DashboardBaseURL = "https://console-api.pinecone.io"
const StagingDashboardBaseURL = "https://staging.console-api.pinecone.io"

func GetDashboardBaseURL() string {
	if config.Staging.Get() {
		return StagingDashboardBaseURL
	} else {
		return DashboardBaseURL
	}
}
