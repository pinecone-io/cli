package auth

import (
	"fmt"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/auth"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/spf13/cobra"
)

func NewCmdAuthStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Show the current authentication status of the Pinecone CLI",
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runAuthStatus(cmd); err != nil {
				log.Error().Err(err).Msg("Error retrieving authentication status")
				exit.Error(pcio.Errorf("error retrieving authentication status: %w", err))
			}
		},
	}
	return cmd
}

func runAuthStatus(cmd *cobra.Command) error {
	token, err := auth.Token(cmd.Context())
	if err != nil { // This should only error on a network request to refresh the token
		log.Error().Err(err).Msg("Error retrieving oauth token")
	}

	apiKey := secrets.GlobalApiKey.Get()
	clientID := secrets.ClientId.Get()
	clientSecret := secrets.ClientSecret.Get()

	authMode := "none"
	switch {
	case token.AccessToken != "" || token.RefreshToken != "":
		authMode = "oauth2-user"
	case clientID != "" && clientSecret != "":
		authMode = "service-account-credentials"
	case apiKey != "":
		authMode = "api-key"
	}

	environment := config.Environment.Get()

	var claims *auth.MyCustomClaims
	if token.AccessToken != "" {
		claims, _ = auth.ParseClaimsUnverified(token)
	}

	expStr := "<none>"
	remaining := ""
	if !token.Expiry.IsZero() {
		expStr = token.Expiry.Format(time.RFC3339)
		remaining = time.Until(token.Expiry).Round(time.Second).String()
	}

	scope := ""
	orgId := ""
	if claims != nil {
		scope = claims.Scope
		orgId = claims.OrgId
	}

	type row struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	rows := []row{
		{Key: "Environment", Value: environment},
		{Key: "Authentication Mode", Value: authMode},
		{Key: "Token Organization ID", Value: orgId},
		{Key: "Token Expiry", Value: expStr},
		{Key: "Token Remaining", Value: remaining},
		{Key: "Token Scope", Value: scope},
	}

	for _, r := range rows {
		fmt.Printf("%-20s %s\n", r.Key, r.Value)
	}

	return nil
}
