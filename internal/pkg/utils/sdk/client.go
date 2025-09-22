package sdk

import (
	"context"
	"crypto/rand"
	"io"

	"github.com/pinecone-io/cli/internal/pkg/utils/auth"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/environment"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"

	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

const (
	CLIAPIKeyName = "pinecone-cli-"
	CLISourceTag  = "pinecone-cli"
)

func NewPineconeClient() *pinecone.Client {
	targetOrgId := state.TargetOrg.Get().Id
	targetProjectId := state.TargetProj.Get().Id
	log.Debug().
		Str("targetOrgId", targetOrgId).
		Str("targetProjectId", targetProjectId).
		Msg("Loading target context")

	ctx := context.Background()
	oauth2Token, err := auth.Token(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving oauth token")
	}
	clientId := secrets.ClientId.Get()
	clientSecret := secrets.ClientSecret.Get()
	globalAPIKey := secrets.GlobalApiKey.Get()

	// If there's a global API key set, it takes priority over user/service account tokens and associated keys
	if secrets.GlobalApiKey.Get() != "" {
		if oauth2Token.AccessToken != "" {
			msg.WarnMsg("You are currently logged in and also have an API key set in your environment and/or local configuration. The API key (which is linked to a specific project) will be used in preference to any user authentication and target context that may be present.\n")
		}

		log.Debug().Msg("Creating client for machine using stored API key")
		return NewClientForAPIKey(secrets.GlobalApiKey.Get())
	}
	log.Debug().Msg("No global API key is stored in configuration, attempting to create a client using user access token")

	// If neither user token or service account credentials are set, we cannot instantiate a client
	if oauth2Token.AccessToken == "" && (clientId == "" && clientSecret == "") && globalAPIKey == "" {
		msg.FailMsg("Please configure user credentials before attempting this operation. Log in with %s, configure a service account with %s, or set an explicit API key with %s.", style.Code("pc login"), style.Code("pc auth configure --client-id --client-secret"), style.Code("pc config set-api-key"))
		exit.ErrorMsg("User credentials are not configured")
	}

	// Lastly, if the user is not targeting a project, we cannot instantiate the client for
	// user or service account scenarios
	if targetProjectId == "" {
		msg.FailMsg("You are logged in, but need to target a project with %s", style.Code("pc target"))
		exit.ErrorMsg("No target project set")
	}

	return NewPineconeClientForProjectById(targetProjectId)
}

func NewPineconeClientForProjectById(projectId string) *pinecone.Client {
	ac := NewPineconeAdminClient()
	ctx := context.Background()

	project, err := ac.Project.Describe(ctx, projectId)
	if err != nil {
		msg.FailMsg("Failed to get project %s: %s", style.Emphasis(projectId), err)
		exit.Error(err)
	}

	// Get the stored ManagedKey for the project, a new key is created if one doesn't exist
	key, err := getCLIAPIKeyForProject(ctx, ac, project)
	if err != nil {
		msg.FailMsg("Failed to retrieve or create an API key for the project %s (ID: %s)", project.Name, project.Id)
		exit.Error(pcio.Errorf("failed to retrieve or create API key for project: %w", err))
	}

	// Header is required for allowing user token to work across data/control plane APIs
	headers := make(map[string]string)
	headers["X-Project-Id"] = projectId

	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey:    key,
		SourceTag: CLISourceTag,
		Host:      getPineconeHostURL(),
		Headers:   headers,
	})
	if err != nil {
		msg.FailMsg("Failed to create Pinecone client: %s", err)
		exit.Error(pcio.Errorf("failed to create Pinecone Client: %w", err))
	}

	return pc
}

func NewClientForAPIKey(apiKey string) *pinecone.Client {
	if apiKey == "" {
		exit.Error(pcio.Errorf("API key not set. Please run %s", style.Code("pc config set-api-key")))
	}

	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey:    apiKey,
		SourceTag: CLISourceTag,
		Host:      getPineconeHostURL(),
	})
	if err != nil {
		exit.Error(err)
	}

	return pc
}

func NewPineconeAdminClient() *pinecone.AdminClient {
	ctx := context.Background()
	oauth2Token, err := auth.Token(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving oauth token")
	}
	clientId := secrets.ClientId.Get()
	clientSecret := secrets.ClientSecret.Get()

	// AdminClient can accept either user token or service account credentials
	// If both are provided, the client will use the user token
	if oauth2Token.AccessToken == "" && (clientId == "" || clientSecret == "") {
		msg.FailMsg("Please login with %s or configure credentials with %s before attempting this operation.", style.Code("pc auth login"), style.Code("pc auth configure"))
		exit.ErrorMsg("User is not logged in")
	}

	sourceTag := CLISourceTag
	ac, err := pinecone.NewAdminClient(pinecone.NewAdminClientParams{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		AccessToken:  oauth2Token.AccessToken,
		SourceTag:    &sourceTag,
		Host:         getPineconeHostURL(),
	})
	if err != nil {
		msg.FailMsg("Failed to create Pinecone admin client: %s", err)
		exit.Error(err)
	}

	return ac
}

func getCLIAPIKeyForProject(ctx context.Context, ac *pinecone.AdminClient, project *pinecone.Project) (string, error) {
	projectAPIKeysMap := secrets.ManagedAPIKeys.Get()
	var managedKey secrets.ManagedKey

	// If we have a managed key stored for the project, use the value
	managedKey, mkExists := projectAPIKeysMap[project.Id]
	if mkExists {
		managedKey = projectAPIKeysMap[project.Id]
		if managedKey.Value != "" {
			return managedKey.Value, nil
		}
	}

	// If we don't have a managed key at this point, we need to create a new one
	newKeyName := generateCLIAPIKeyName()
	newKey, err := ac.APIKey.Create(ctx, project.Id, &pinecone.CreateAPIKeyParams{
		Name: newKeyName,
	})
	if err != nil {
		msg.FailMsg("Failed to create a CLI managed API key for project %s: %s", style.Emphasis(project.Name), err)
		return "", pcio.Errorf("failed to create a CLI managed API key for project: %w", err)
	}

	managedKey = secrets.ManagedKey{
		ProjectId:      project.Id,
		OrganizationId: project.OrganizationId,
		Value:          newKey.Value,
		Name:           newKeyName,
		Origin:         secrets.OriginCLICreated,
	}

	// Add the new ManagedKey to the map
	projectAPIKeysMap[project.Id] = managedKey
	secrets.ManagedAPIKeys.Set(projectAPIKeysMap)

	return managedKey.Value, nil
}

func getPineconeHostURL() string {
	env := config.Environment.Get()
	connectionConfig, err := environment.GetEnvConfig(env)
	if err != nil { // If there's an error resolving the environment, default to production host
		return environment.Prod.PineconeGCPURL
	}
	return connectionConfig.PineconeGCPURL
}

func randStringFromCharset(length int) (string, error) {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	idxMax := 256 - (256 % len(charset))
	out := make([]byte, length)
	var b [1]byte
	for i := 0; i < length; {
		if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
			return "", err
		}
		if int(b[0]) >= idxMax {
			continue // reject 252-255 to keep uniform distribution
		}
		out[i] = charset[int(b[0])%len(charset)]
		i++
	}
	return string(out), nil
}

func generateCLIAPIKeyName() string {
	const suffixLength = 6
	s, err := randStringFromCharset(suffixLength)
	if err != nil {
		return CLIAPIKeyName + "000000" // fallback if randomization errors
	}
	return CLIAPIKeyName + s
}
