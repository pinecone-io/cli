package sdk

import (
	"context"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/environment"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	pc_oauth2 "github.com/pinecone-io/cli/internal/pkg/utils/oauth2"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"

	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

func newClientParams(key string) pinecone.NewClientParams {
	env := config.Environment.Get()

	var clientControllerHostUrl string
	switch env {
	case "production":
		clientControllerHostUrl = environment.Prod.IndexControlPlaneUrl
	case "staging":
		clientControllerHostUrl = environment.Staging.IndexControlPlaneUrl
	default:
		exit.Error(pcio.Errorf("invalid environment: %s", env))
		return pinecone.NewClientParams{}
	}

	return pinecone.NewClientParams{
		ApiKey:    key,
		SourceTag: "pinecone-cli",
		Host:      clientControllerHostUrl,
	}
}

func newClientForUserFromTarget() *pinecone.Client {
	targetOrgId := state.TargetOrg.Get().Id
	targetProjectId := state.TargetProj.Get().Id
	log.Debug().
		Str("targetOrgId", targetOrgId).
		Str("targetProjectId", targetProjectId).
		Msg("Loading target context")

	oauth2Token := secrets.OAuth2Token.Get()

	if secrets.ApiKey.Get() != "" {
		if oauth2Token.AccessToken != "" {
			msg.WarnMsg("You are currently logged in and also have an API key set in your environment and/or local configuration. The API key (which is linked to a specific project) will be used in preference to any user authentication and target context that may be present.\n")
		}

		log.Debug().Msg("Creating client for machine using stored API key")
		return NewClientForMachine(secrets.ApiKey.Get())
	}

	log.Debug().Msg("No API key is stored in configuration, so attempting to create a client using user access token")

	if oauth2Token.AccessToken == "" {
		msg.FailMsg("Please set an API key with %s or login with %s before attempting this operation.", style.Code("pinecone config set-api-key"), style.Code("pinecone login"))
		exit.ErrorMsg("User is not logged in")
	}

	if targetOrgId == "" || targetProjectId == "" {
		msg.FailMsg("You are logged in, but need to target a project with %s", style.Code("pinecone target"))
		exit.ErrorMsg("No target organization set")
	}

	return NewPineconeClientForUser(targetProjectId)
}

func NewPineconeClientForUser(projectId string) *pinecone.Client {
	env := config.Environment.Get()
	connectionConfig, err := environment.GetEnvConfig(env)
	if err != nil {
		msg.FailMsg("Failed to get connection configuration for environment %s: %s", env, err)
		exit.Error(err)
	}

	headers := make(map[string]string)
	headers["X-Project-Id"] = projectId

	ctx := context.Background()
	targetOrgId := state.TargetOrg.Get().Id
	restClient, err := pc_oauth2.GetHttpClient(ctx, &targetOrgId)
	if err != nil {
		msg.FailMsg("Failed to create OAuth2 client: %s", err)
		exit.Error(err)
	}

	pc, err := pinecone.NewClientBase(pinecone.NewClientBaseParams{
		Host:       connectionConfig.IndexControlPlaneUrl,
		Headers:    headers,
		RestClient: restClient,
	})
	if err != nil {
		msg.FailMsg("Failed to create Pinecone client: %s", err)
		exit.Error(err)
	}

	return pc
}

func NewClientForMachine(apiKey string) *pinecone.Client {
	if apiKey == "" {
		exit.Error(pcio.Errorf("API key not set. Please run %s", style.Code("pinecone config set-api-key")))
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
		msg.FailMsg("API key not set. Please run %s or %s", style.Code("pinecone login"), style.Code("pinecone config set-api-key"))
		exit.Error(pcio.Errorf("API key not set."))
	}

	pc, err := pinecone.NewClient(newClientParams(key))
	if err != nil {
		msg.FailMsg("Failed to create Pinecone client: %s", err)
		exit.Error(err)
	}

	return pc
}
