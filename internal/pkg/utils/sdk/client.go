package sdk

import (
	"context"

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

const (
	cliApiKeyName = "pinecone-cli-api-key"
)

func newClientParams(apiKey string) pinecone.NewClientParams {
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
		ApiKey:    apiKey,
		SourceTag: "pinecone-cli",
		Host:      clientControllerHostUrl,
	}
}

func newAdminClientParams(clientId string, clientSecret string, accessToken string) pinecone.NewAdminClientParams {
	env := config.Environment.Get()

	var clientControllerHostUrl string
	switch env {
	case "production":
		clientControllerHostUrl = environment.Prod.IndexControlPlaneUrl
	case "staging":
		clientControllerHostUrl = environment.Staging.IndexControlPlaneUrl
	default:
		exit.Error(pcio.Errorf("invalid environment: %s", env))
		return pinecone.NewAdminClientParams{}
	}
	return pinecone.NewAdminClientParams{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
		Host:         clientControllerHostUrl,
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

	log.Debug().Msg("No API key is stored in configuration, attempting to create a client using user access token")

	if oauth2Token.AccessToken == "" {
		msg.FailMsg("Please set an API key with %s or login with %s before attempting this operation.", style.Code("pc config set-api-key"), style.Code("pc login"))
		exit.ErrorMsg("User is not logged in")
	}

	if targetOrgId == "" || targetProjectId == "" {
		msg.FailMsg("You are logged in, but need to target a project with %s", style.Code("pc target"))
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
		exit.Error(pcio.Errorf("API key not set. Please run %s", style.Code("pc config set-api-key")))
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

func NewPineconeAdminClient() *pinecone.AdminClient {
	oauth2Token := secrets.OAuth2Token.Get()
	clientId := secrets.ClientId.Get()
	clientSecret := secrets.ClientSecret.Get()

	if oauth2Token.AccessToken == "" && (clientId == "" || clientSecret == "") {
		msg.FailMsg("Please login with %s or configure credentials with %s before attempting this operation.", style.Code("pc auth login"), style.Code("pc auth configure"))
		exit.ErrorMsg("User is not logged in")
	}

	ac, err := pinecone.NewAdminClient(newAdminClientParams(clientId, clientSecret, oauth2Token.AccessToken))
	if err != nil {
		msg.FailMsg("Failed to create Pinecone admin client: %s", err)
		exit.Error(err)
	}

	return ac
}

func NewPineconeClientForProjectById(projectId string) *pinecone.Client {
	ac := NewPineconeAdminClient()
	ctx := context.Background()

	project, err := ac.Project.Describe(ctx, projectId)
	if err != nil {
		msg.FailMsg("Failed to get project %s: %s", style.Emphasis(projectId), err)
		exit.Error(err)
	}

	key, err := getCLIAPIKeyForProject(ctx, ac, project)
	if err != nil {
		msg.FailMsg("Failed to retrieve or create an API key for the project %s (ID: %s)", project.Name, project.Id)
		exit.Error(pcio.Errorf("failed to retrieve or create API key for project: %w", err))
	}

	pc, err := pinecone.NewClient(newClientParams(key))
	if err != nil {
		msg.FailMsg("Failed to create Pinecone client: %s", err)
		exit.Error(pcio.Errorf("failed to create Pinecone Client: %w", err))
	}

	return pc
}

func getCLIAPIKeyForProject(ctx context.Context, ac *pinecone.AdminClient, project *pinecone.Project) (string, error) {
	apiKeys, err := ac.APIKey.List(ctx, project.Id)
	if err != nil {
		msg.FailMsg("Failed to get API keys for project %s: %s", style.Emphasis(project.Name), err)
		exit.Error(err)
	}

	projectAPIKeysMap := secrets.ProjectAPIKeys.Get()

	var keyValue string
	var existingProjectAPIKey *pinecone.APIKey
	projectHasCLIAPIKey := false
	if len(apiKeys) > 0 {
		for _, key := range apiKeys {
			if key.Name == cliApiKeyName {
				projectHasCLIAPIKey = true
				existingProjectAPIKey = key
				break
			}
		}

		// if the project has a CLI API key currently via listing API keys
		if projectHasCLIAPIKey {
			// if the project has a CLI API key stored in secrets state, use it
			if projectAPIKeysMap[project.Id] != "" {
				keyValue = projectAPIKeysMap[project.Id]

				return keyValue, nil
			}

			// if the project does not have an associated CLI API key stored in state, delete the existing key
			// and create a new one below
			if projectAPIKeysMap[project.Id] == "" {
				err := ac.APIKey.Delete(ctx, existingProjectAPIKey.Id)
				if err != nil {
					msg.FailMsg("Failed to delete API key for project %s: %s", style.Emphasis(project.Name), err)
					exit.Error(err)
				}
			}
		}
		// if keyValue is still empty, we know we need to create a new key and persist in secrets state below
	}

	if keyValue == "" {
		// create a new CLI API key
		newKey, err := ac.APIKey.Create(ctx, project.Id, &pinecone.CreateAPIKeyParams{
			Name: cliApiKeyName,
		})
		if err != nil {
			msg.FailMsg("Failed to create CLI API key for project %s: %s", style.Emphasis(project.Name), err)
			return "", pcio.Errorf("failed to create a CLI API key for project: %w", err)
		}

		keyValue = newKey.Value

		// persist the new key in secrets state
		projectAPIKeysMap[project.Id] = keyValue
		secrets.ProjectAPIKeys.Set(&projectAPIKeysMap)
	}

	return keyValue, nil
}
