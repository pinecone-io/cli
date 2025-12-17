//go:build e2e

package e2e

import (
	"github.com/pinecone-io/cli/test/e2e/helpers"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func (s *ServiceAccountSuite) TestAPIKeyLifecycle() {
	// Requires admin client to manage API keys
	_, _ = helpers.RequireServiceAccount(s.T())

	projID := helpers.ProjectID()
	if projID == "" {
		s.T().Skip("PC_E2E_PROJECT_ID not set; skipping api-key lifecycle test")
	}

	name := helpers.RandomName("e2e-key")

	var create pinecone.APIKeyWithSecret
	_, err := s.cli.RunJSONCtx(s.ctx, &create, "api-key", "create", "--id", projID, "--name", name)
	s.Require().NoError(err, "api-key create failed")
	s.Require().NotEmpty(create.Key.Id, "expected created key id")
	s.Require().NotEmpty(create.Value, "expected created key value")

	var desc pinecone.APIKey
	_, err = s.cli.RunJSONCtx(s.ctx, &desc, "api-key", "describe", "--id", create.Key.Id)
	s.Require().NoError(err, "api-key describe failed")
	s.Require().Equal(create.Key.Id, desc.Id, "describe id mismatch")

	s.T().Cleanup(func() {
		_, _, err := s.cli.RunCtx(s.ctx, "api-key", "delete", "--id", desc.Id, "--skip-confirmation")
		s.Require().NoError(err, "api-key delete failed")
	})
}
