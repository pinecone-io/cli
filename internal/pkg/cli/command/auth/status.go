package auth

import (
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type authStatusCmdOptions struct {
	json bool
}

func NewCmdAuthStatus() *cobra.Command {
	options := authStatusCmdOptions{}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the current authentication configuration for the Pinecone CLI",
		Example: help.Examples(`
			pc auth status --json
		`),
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runAuthStatus(cmd, options); err != nil {
				msg.FailMsg("Error retrieving authentication status: %s", err)
				exit.Error().Err(err).Msg("Error retrieving authentication status")
			}
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func runAuthStatus(cmd *cobra.Command, options authStatusCmdOptions) error {
	token, err := oauth.Token(cmd.Context())
	if err != nil { // This should only error on a network request to refresh the token
		return err
	}

	authMode := string(state.AuthedUser.Get().AuthContext)
	orgName := state.TargetOrg.Get().Name
	projName := state.TargetProj.Get().Name
	environment := config.Environment.Get()

	// Default API Key
	defaultAPIKey := secrets.DefaultAPIKey.Get()

	// Service Account
	clientId := secrets.ClientId.Get()
	clientSecret := secrets.ClientSecret.Get()

	// Extract token information
	var claims *oauth.MyCustomClaims
	expStr := ""
	remaining := ""
	scope := ""
	orgId := ""

	if token != nil {
		if token.AccessToken != "" {
			claims, _ = oauth.ParseClaimsUnverified(token)
		}

		if !token.Expiry.IsZero() {
			expStr = token.Expiry.Format(time.RFC3339)
			remaining = time.Until(token.Expiry).Round(time.Second).String()
		}

		if claims != nil {
			scope = claims.Scope
			orgId = claims.OrgId
		}
	}
	authStatus := presenters.AuthStatus{
		AuthMode:            authMode,
		OrganizationName:    orgName,
		ProjectName:         projName,
		Token:               token,
		DefaultAPIKey:       presenters.MaskHeadTail(defaultAPIKey, 4, 4),
		ClientID:            clientId,
		ClientSecret:        presenters.MaskHeadTail(clientSecret, 4, 4),
		TokenExpiry:         expStr,
		TokenTimeRemaining:  remaining,
		TokenScope:          scope,
		TokenOrganizationID: orgId,
		Environment:         environment,
	}

	if options.json {
		json := text.IndentJSON(authStatus)
		pcio.Println(json)
		return nil
	}

	presenters.PrintAuthStatus(authStatus)

	return nil
}
