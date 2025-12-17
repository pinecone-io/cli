//go:build e2e

package e2e

import (
	"github.com/pinecone-io/cli/test/e2e/helpers"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

// Service account tests
func (s *ServiceAccountSuite) TestNamespaceLifecycle() {
	if s.indexName == "" {
		s.T().Skip("no index available for namespace tests")
	}

	name := helpers.RandomName("e2e-namespace")
	var namespace pinecone.NamespaceDescription
	_, err := s.cli.RunJSONCtx(s.ctx, &namespace, "index", "namespace", "create", "--index-name", s.indexName, "--name", name)
	s.Require().NoError(err, "namespace create failed")
	s.Require().Equal(name, namespace.Name, "namespace name mismatch")

	var describe pinecone.NamespaceDescription
	_, err = s.cli.RunJSONCtx(s.ctx, &describe, "index", "namespace", "describe", "--index-name", s.indexName, "--name", name)
	s.Require().NoError(err, "namespace describe failed")
	s.Require().Equal(name, describe.Name, "namespace name mismatch")

	var list pinecone.ListNamespacesResponse
	_, err = s.cli.RunJSONCtx(s.ctx, &list, "index", "namespace", "list", "--index-name", s.indexName)
	s.Require().NoError(err, "namespace list failed")
	s.Require().Equal(1, len(list.Namespaces), "expected one namespace in list")
	s.Require().Equal(name, list.Namespaces[0].Name, "namespace name mismatch")

	_, _, err = s.cli.RunCtx(s.ctx, "index", "namespace", "delete", "--index-name", s.indexName, "--name", name)
	s.Require().NoError(err, "namespace delete failed")
}

// API key tests
func (a *APIKeySuite) TestNamespaceLifecycle() {
	if a.indexName == "" {
		a.T().Skip("no index available for namespace tests")
	}

	name := helpers.RandomName("e2e-namespace")
	var namespace pinecone.NamespaceDescription
	_, err := a.cli.RunJSONCtx(a.ctx, &namespace, "index", "namespace", "create", "--index-name", a.indexName, "--name", name)
	a.Require().NoError(err, "namespace create failed")
	a.Require().Equal(name, namespace.Name, "namespace name mismatch")

	var describe pinecone.NamespaceDescription
	_, err = a.cli.RunJSONCtx(a.ctx, &describe, "index", "namespace", "describe", "--index-name", a.indexName, "--name", name)
	a.Require().NoError(err, "namespace describe failed")
	a.Require().Equal(name, describe.Name, "namespace name mismatch")

	var list pinecone.ListNamespacesResponse
	_, err = a.cli.RunJSONCtx(a.ctx, &list, "index", "namespace", "list", "--index-name", a.indexName)
	a.Require().NoError(err, "namespace list failed")
	a.Require().Equal(1, len(list.Namespaces), "expected one namespace in list")
	a.Require().Equal(name, list.Namespaces[0].Name, "namespace name mismatch")

	_, _, err = a.cli.RunCtx(a.ctx, "index", "namespace", "delete", "--index-name", a.indexName, "--name", name)
	a.Require().NoError(err, "namespace delete failed")
}
