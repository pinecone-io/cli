package client

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

func newClientParams(key string) pinecone.NewClientParams {
	return pinecone.NewClientParams{
		ApiKey:    key,
		SourceTag: "pinecone-cli",
	}
}

func newClientForUser() *pinecone.Client {
	target := state.GetTargetContext()

	if target.Org == "" || target.Project == "" {
		fmt.Println("Target context is currently:")
		fmt.Println()
		presenters.PrintTargetContext(target)
		fmt.Println()
		exit.Error(fmt.Errorf("The target organization and project must both be set. Please run %s", style.Code("pinecone target")))
	}

	orgs, err := dashboard.GetOrganizations(secrets.AccessToken.Get())
	if err != nil {
		exit.Error(err)
	}

	var project dashboard.Project
	for _, org := range orgs.Organizations {
		if org.Name == target.Org {
			for _, proj := range org.Projects {
				if proj.Name == target.Project {
					project = proj
					break
				}
			}
		}
	}

	keyResponse, err2 := dashboard.GetApiKeys(project, secrets.AccessToken.Get())
	if err2 != nil {
		exit.Error(err2)
	}

	var key string
	if len(keyResponse.Keys) > 0 {
		key = keyResponse.Keys[0].Value
	} else {
		exit.Error(fmt.Errorf("No API keys found for project %s", style.Code(target.Project)))
	}

	if key == "" {
		exit.Error(fmt.Errorf("API key not set. Please run %s or %s", style.Code("pinecone auth login"), style.Code("pinecone auth set-api-key")))
	}

	pc, err := pinecone.NewClient(newClientParams(key))
	if err != nil {
		exit.Error(err)
	}

	return pc
}

func newClientForMachine() *pinecone.Client {
	key := secrets.ApiKey.Get()
	if key == "" {
		exit.Error(fmt.Errorf("API key not set. Please run %s", style.Code("pinecone auth set-api-key")))
	}

	pc, err := pinecone.NewClient(newClientParams(key))
	if err != nil {
		exit.Error(err)
	}

	return pc
}

func NewPineconeClient() *pinecone.Client {
	key := secrets.ApiKey.Get()
	if key == "" {
		return newClientForUser()
	} else {
		return newClientForMachine()
	}
}
