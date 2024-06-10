package sdk

import (
	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/environment"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

func newClientParams(key string) pinecone.NewClientParams {
	env := config.Environment.Get()

	var clientControllerHostUrl string
	if env == "production" {
		clientControllerHostUrl = environment.Prod.IndexControlPlaneUrl
	} else {
		clientControllerHostUrl = environment.Staging.IndexControlPlaneUrl
	}

	return pinecone.NewClientParams{
		ApiKey:    key,
		SourceTag: "pinecone-cli",
		Host:      clientControllerHostUrl,
	}
}

func newClientForUserFromTarget() *pinecone.Client {
	targetOrgName := state.TargetOrg.Get().Name
	targetOrgId := state.TargetOrg.Get().Id
	targetProjectName := state.TargetProj.Get().Name
	targetProjectId := state.TargetProj.Get().Id

	apiKey := secrets.ApiKey.Get()

	if targetOrgId == "" || targetProjectId == "" {

		if apiKey != "" {
			return NewClientForMachine(apiKey)
		}

		msg.FailMsg("Please run %s to set a target context", style.Code("pinecone target"))
		pcio.Println()
		pcio.Println("Target context is currently:")
		pcio.Println()
		presenters.PrintTargetContext(state.GetTargetContext())
		pcio.Println()
		exit.ErrorMsg(pcio.Sprintf("The target organization and project must both be set. Please run %s", style.Code("pinecone target")))
	}

	return NewPineconeClientForUserProjectByName(targetOrgName, targetProjectName)
}

func NewPineconeClientForUserProjectByName(orgName string, projectName string) *pinecone.Client {
	project, err := dashboard.GetProjectByName(orgName, projectName)
	if err != nil {
		msg.FailMsg("Failed to get project %s: %s", style.Emphasis(projectName), err)
		exit.Error(err)
	}

	keyResponse, err2 := dashboard.GetApiKeys(*project)
	if err2 != nil {
		msg.FailMsg("Failed to get API keys for project %s: %s", style.Emphasis(project.Name), err2)
		exit.Error(err2)
	}

	var key string
	if len(keyResponse.Keys) > 0 {
		key = keyResponse.Keys[0].Value
	} else {
		log.Error().Str("projectName", projectName).Msg("No API keys found for project")
		msg.FailMsg("No API keys found for project %s", style.Code(projectName))
		exit.ErrorMsg(pcio.Sprintf("No API keys found for project %s", style.Emphasis(projectName)))
	}

	if key == "" {
		msg.FailMsg("API key not set. Please run %s", style.Code("pinecone auth login"))
		exit.Error(pcio.Errorf("API key not set."))
	}

	pc, err := pinecone.NewClient(newClientParams(key))
	if err != nil {
		msg.FailMsg("Failed to create Pinecone client: %s", err)
		exit.Error(err)
	}

	return pc
}

func NewClientForMachine(apiKey string) *pinecone.Client {
	if apiKey == "" {
		exit.Error(pcio.Errorf("API key not set. Please run %s", style.Code("pinecone auth set-api-key")))
	}

	pc, err := pinecone.NewClient(newClientParams(apiKey))
	if err != nil {
		exit.Error(err)
	}

	return pc
}

func NewPineconeClient() *pinecone.Client {
	return newClientForUserFromTarget()
}

func NewPineconeClientForProjectById(orgId string, projectId string) *pinecone.Client {
	project, err := dashboard.GetProjectById(orgId, projectId)
	if err != nil {
		msg.FailMsg("Failed to get project %s: %s", style.Emphasis(projectId), err)
		exit.Error(err)
	}

	keyResponse, err2 := dashboard.GetApiKeys(*project)
	if err2 != nil {
		msg.FailMsg("Failed to get API keys for project %s: %s", style.Emphasis(project.Name), err2)
		exit.Error(err2)
	}

	var key string
	if len(keyResponse.Keys) > 0 {
		key = keyResponse.Keys[0].Value
	} else {
		log.Error().Str("projectId", projectId).Msg("No API keys found for project")
		msg.FailMsg("No API keys found for project id %s", style.Code(projectId))
		exit.ErrorMsg(pcio.Sprintf("No API keys found for project %s", style.Emphasis(projectId)))
	}

	if key == "" {
		msg.FailMsg("API key not set. Please run %s", style.Code("pinecone auth login"))
		exit.Error(pcio.Errorf("API key not set."))
	}

	pc, err := pinecone.NewClient(newClientParams(key))
	if err != nil {
		msg.FailMsg("Failed to create Pinecone client: %s", err)
		exit.Error(err)
	}

	return pc
}
