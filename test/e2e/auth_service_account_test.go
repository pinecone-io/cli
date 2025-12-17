//go:build e2e

package e2e

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/test/e2e/helpers"
)

func (s *ServiceAccountSuite) TestAuthServiceAccountConfigureAndStatus() {
	// Requires service account credentials for calling pc auth configure
	clientID, clientSecret := helpers.RequireServiceAccount(s.T())

	var context state.TargetContext
	_, err := s.cli.RunJSONCtx(s.ctx, &context,
		"auth", "configure",
		"--client-id", clientID,
		"--client-secret", clientSecret,
		"--project-id", helpers.ProjectID(),
		"--prompt-if-missing=false",
	)
	s.Require().NoError(err, "auth configure/status failed")

	s.Require().NotEmpty(context.Organization.Id, "expected organization id after configure")
	s.Require().NotEmpty(context.Project.Id, "expected project id after configure")
}
