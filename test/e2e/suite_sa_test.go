//go:build e2e

package e2e

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/pinecone-io/cli/test/e2e/helpers"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

type ServiceAccountSuite struct {
	suite.Suite
	cli        *helpers.CLI
	ctx        context.Context
	indexName  string
	index      pinecone.Index
	createdIdx bool
}

func TestServiceAccountSuite(t *testing.T) {
	suite.Run(t, new(ServiceAccountSuite))
}

func (s *ServiceAccountSuite) SetupSuite() {
	helpers.RequireE2E(s.T())
	_, _ = helpers.RequireServiceAccount(s.T())

	s.cli = helpers.NewCLI(s.T())
	s.ctx = context.Background()

	name := helpers.RandomName("e2e-srvless-sa")
	args := []string{
		"index", "create",
		"--name", name,
		"--cloud", helpers.Cloud(),
		"--region", helpers.Region(),
		"--dimension", strconv.Itoa(helpers.Dimension()),
		"--metric", "cosine",
	}
	_, err := s.cli.RunJSONCtx(s.ctx, &s.index, args...)
	s.Require().NoError(err, "index create failed (sa)")
	s.indexName = name
	s.createdIdx = true

	err = helpers.WaitForIndexReady(s.cli, name, 5*time.Minute)
	s.Require().NoError(err, "index not ready (sa)")

}

func (s *ServiceAccountSuite) TearDownSuite() {
	if !s.createdIdx || s.indexName == "" || s.cli == nil {
		return
	}
	_, _, err := s.cli.RunCtx(s.ctx, "index", "delete", "--name", s.indexName)
	s.Require().NoError(err, "index delete failed (sa)")
}
