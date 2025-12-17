//go:build e2e

package e2e

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/pinecone-io/cli/test/e2e/helpers"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
)

type APIKeySuite struct {
	suite.Suite
	cli          *helpers.CLI
	ctx          context.Context
	indexName    string
	createdIndex bool
	tempHome     string
	prevHome     string
}

func TestAPIKeySuite(t *testing.T) {
	suite.Run(t, new(APIKeySuite))
}

func (a *APIKeySuite) SetupSuite() {
	helpers.RequireE2E(a.T())
	_ = helpers.RequireAPIKey(a.T())

	// Fresh home so no service-account credentials leak into the API-key suite.
	a.prevHome = os.Getenv("PC_E2E_HOME")
	tempHome, _ := os.MkdirTemp("", "pc-e2e-home-api-*")
	a.tempHome = tempHome
	_ = os.Setenv("PC_E2E_HOME", tempHome)
	a.T().Cleanup(func() {
		if a.prevHome == "" {
			_ = os.Unsetenv("PC_E2E_HOME")
		} else {
			_ = os.Setenv("PC_E2E_HOME", a.prevHome)
		}
		_ = os.RemoveAll(tempHome)
	})

	a.cli = helpers.NewCLI(a.T())
	a.ctx = context.Background()

	name := helpers.RandomName("e2e-api")
	args := []string{
		"index", "create",
		"--name", name,
		"--cloud", helpers.Cloud(),
		"--region", helpers.Region(),
		"--dimension", strconv.Itoa(helpers.Dimension()),
		"--metric", "cosine",
	}
	var idx pinecone.Index
	_, err := a.cli.RunJSONCtx(a.ctx, &idx, args...)
	a.Require().NoError(err, "index create failed (api key)")
	a.indexName = name
	a.createdIndex = true

	if err := helpers.WaitForIndexReady(a.cli, name, 5*time.Minute); err != nil {
		a.T().Fatalf("index not ready (api key): %v", err)
	}
}

func (a *APIKeySuite) TearDownSuite() {
	if a.createdIndex && a.indexName != "" && a.cli != nil {
		_, _, _ = a.cli.RunCtx(a.ctx, "index", "delete", "--name", a.indexName)
	}
	if a.tempHome != "" {
		_ = os.RemoveAll(a.tempHome)
	}
	// Restore PC_E2E_HOME in case any additional suites are run.
	if a.prevHome == "" {
		_ = os.Unsetenv("PC_E2E_HOME")
	} else {
		_ = os.Setenv("PC_E2E_HOME", a.prevHome)
	}
}
