package config

import (
	conf "github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewSetEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-environment <production|staging>",
		Short: "Configure the environment (production or staging)",
		Example: help.Examples(`
			pc config set-environment production
			pc config set-environment staging
		`),
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				msg.FailMsg("Please provide a value for environment. Accepted values are %s, %s", style.Emphasis("production"), style.Emphasis("staging"))
				exit.ErrorMsg("No value provided for environment")
			}
			envArg := args[0]

			var settingValue string
			switch envArg {
			case "staging":
				settingValue = "staging"
			case "production", "prod":
				settingValue = "production"
			default:
				msg.FailMsg("Invalid environment. Please use %s or %s.", style.Emphasis("staging"), style.Emphasis("production"))
				exit.ErrorMsg("Invalid environment " + envArg)
			}

			oldConfig := conf.Environment.Get()
			if oldConfig == settingValue {
				msg.InfoMsg("Environment is already set to %s", style.Emphasis(settingValue))
				return
			}

			conf.Environment.Set(settingValue)
			msg.SuccessMsg("Config property %s updated to %s", style.Emphasis("environment"), style.Emphasis(settingValue))

			token, err := oauth.Token(cmd.Context())
			if err != nil {
				log.Error().Err(err).Msg("Error retrieving oauth token")
				msg.FailMsg("Error retrieving oauth token: %s", err)
				exit.Error(pcio.Errorf("error retrieving oauth token: %w", err))
			}
			if token.AccessToken != "" || token.RefreshToken != "" {
				oauth.Logout()
				msg.InfoMsg("You have been logged out; to login again, run %s", style.Code("pc login"))
			} else {
				msg.InfoMsg("To login, run %s", style.Code("pc login"))
			}

			if secrets.GlobalApiKey.Get() != "" {
				secrets.GlobalApiKey.Clear()
				msg.InfoMsg("API key cleared; to set a new API key, run %s", style.Code("pc config set-api-key"))
			} else {
				msg.InfoMsg("To set a new API key, run %s", style.Code("pc config set-api-key"))
			}

			if state.TargetOrg.Get().Name != "" || state.TargetProj.Get().Name != "" {
				state.TargetOrg.Clear()
				state.TargetProj.Clear()
				msg.InfoMsg("Target organization and project cleared; to set a new target, run %s", style.Code("pc target -o myorg -p myproj"))
			}
		},
	}

	return cmd
}
