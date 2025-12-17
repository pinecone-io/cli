//go:build e2e

package e2e

import (
	"github.com/pinecone-io/cli/test/e2e/helpers"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func (s *ServiceAccountSuite) TestProjectList() {
	// Requires admin client
	_, _ = helpers.RequireServiceAccount(s.T())

	var projects []pinecone.Project
	_, err := s.cli.RunJSONCtx(s.ctx, &projects, "project", "list")
	s.Require().NoError(err, "project list failed")
	s.Require().NotEmpty(projects, "expected at least one project")
}
