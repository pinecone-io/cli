//go:build e2e

package e2e

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/test/e2e/helpers"
)

func (s *ServiceAccountSuite) TestTargetSetAndShow() {
	// admin needed to resolve org/project by ID
	_, _ = helpers.RequireServiceAccount(s.T())

	orgID := helpers.OrgID()
	projID := helpers.ProjectID()
	if orgID == "" || projID == "" {
		s.T().Skip("PC_E2E_ORG_ID or PC_E2E_PROJECT_ID not set; skipping target test")
	}

	_, _, err := s.cli.RunCtx(s.ctx, "target", "--org-id", orgID, "--project-id", projID)
	s.Require().NoError(err, "target set failed")

	var tc state.TargetContext
	_, err = s.cli.RunJSONCtx(s.ctx, &tc, "target", "--show")
	s.Require().NoError(err, "target --show failed")
	s.Require().Equal(orgID, tc.Organization.Id, "organization id mismatch")
	s.Require().Equal(projID, tc.Project.Id, "project id mismatch")
}
