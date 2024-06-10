package config

import (
	conf "github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewSetEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-environment <production|staging>",
		Short: "Configure the environment (production or staging)",
		Run: func(cmd *cobra.Command, args []string) {
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

			if secrets.OAuth2Token.Get().AccessToken != "" {
				secrets.OAuth2Token.Clear()
				msg.InfoMsg("You have been logged out; to login again, run %s", style.Code("pinecone login"))
			} else {
				msg.InfoMsg("To login, run %s", style.Code("pinecone login"))
			}

			if secrets.ApiKey.Get() != "" {
				secrets.ApiKey.Clear()
				msg.InfoMsg("API key cleared; to set a new API key, run %s", style.Code("pinecone config set-api-key"))
			} else {
				msg.InfoMsg("To set a new API key, run %s", style.Code("pinecone config set-api-key"))
			}

			if (state.TargetOrg.Get().Name != "" || state.TargetProj.Get().Name != "") && state.TargetKm.Get().Name != "" {
				state.TargetOrg.Clear()
				state.TargetProj.Clear()
				msg.InfoMsg("Target organization and project cleared; to set a new target, run %s", style.Code("pinecone target -o myorg -p myproj"))
			}

			if state.TargetKm.Get().Name != "" {
				state.TargetKm.Clear()
				msg.InfoMsg("Target knowledge model cleared; to set a new target model, run %s", style.Code("pinecone km target -m mymodel"))
			}
		},
	}

	return cmd
}