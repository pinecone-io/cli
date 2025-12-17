//go:build e2e

package e2e

import (
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

func (s *ServiceAccountSuite) TestIndexServerless_ServiceAccountDescribeAndList() {
	var desc pinecone.Index
	_, err := s.cli.RunJSONCtx(s.ctx, &desc, "index", "describe", "--name", s.indexName)
	s.Require().NoError(err, "index describe failed")
	s.Require().Equal(s.indexName, desc.Name, "describe name mismatch")

	var list []pinecone.Index
	_, err = s.cli.RunJSONCtx(s.ctx, &list, "index", "list")
	s.Require().NoError(err, "index list failed")
	s.Require().NotEmpty(list, "expected at least one index in list")

	found := false
	for _, idx := range list {
		if idx.Name == s.indexName {
			found = true
			break
		}
	}
	s.Require().True(found, "shared index not found in list output")
}

func (a *APIKeySuite) TestIndexServerless_APIKeyDescribeAndList() {
	if a.indexName == "" {
		a.T().Skip("no shared index available for API key suite")
	}

	var desc pinecone.Index
	_, err := a.cli.RunJSONCtx(a.ctx, &desc, "index", "describe", "--name", a.indexName)
	a.Require().NoError(err, "index describe failed (api key)")
	a.Require().Equal(a.indexName, desc.Name, "describe name mismatch (api key)")

	var list []pinecone.Index
	_, err = a.cli.RunJSONCtx(a.ctx, &list, "index", "list")
	a.Require().NoError(err, "index list failed (api key)")
	a.Require().NotEmpty(list, "expected at least one index in list (api key)")

	found := false
	for _, idx := range list {
		if idx.Name == a.indexName {
			found = true
			break
		}
	}
	a.Require().True(found, "shared index not found in list output (api key)")
}
