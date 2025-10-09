package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/spf13/cobra"
)

type clearCmdOptions struct {
	serviceAccount bool
	defaultAPIKey  bool
}

func NewClearCmd() *cobra.Command {
	options := clearCmdOptions{}

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear a service account (client ID and secret) or API key from local storage",
		Example: help.Examples(`
		    # Clear configured service account credentials
		    pc auth clear --service-account

		    # Clear configured default API key
		    pc auth clear --api-key

			# Clear both configured service account credentials and default API key
			pc auth clear --service-account --api-key
		`),
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			if !options.serviceAccount && !options.defaultAPIKey {
				msg.FailMsg("Please specify either --service-account or --global-api-key")
				exit.ErrorMsg("No option specified")
			}

			if options.serviceAccount {
				secrets.ClientId.Clear()
				secrets.ClientSecret.Clear()
				msg.SuccessMsg("Service account (client ID and secret) cleared from local storage")
			}

			if options.defaultAPIKey {
				secrets.DefaultAPIKey.Clear()
				msg.SuccessMsg("Default API key cleared")
			}

			// After clearing things, we need to resolve whether the user is still authenticated
			if secrets.DefaultAPIKey.Get() != "" {
				state.AuthedUser.Update(func(u *state.TargetUser) {
					u.AuthContext = state.AuthDefaultAPIKey
				})
			} else if secrets.ClientId.Get() != "" && secrets.ClientSecret.Get() != "" {
				state.AuthedUser.Update(func(u *state.TargetUser) {
					u.AuthContext = state.AuthServiceAccount
				})
			} else if secrets.GetOAuth2Token().AccessToken != "" {
				state.AuthedUser.Update(func(u *state.TargetUser) {
					u.AuthContext = state.AuthUserToken
				})
			} else {
				state.AuthedUser.Update(func(u *state.TargetUser) {
					u.AuthContext = state.AuthNone
				})
			}
		},
	}

	cmd.Flags().BoolVar(&options.serviceAccount, "service-account", false, "Clear the configured service account (client ID and secret) from local storage")
	cmd.Flags().BoolVar(&options.defaultAPIKey, "api-key", false, "Clear the default API key from local storage")

	return cmd
}
