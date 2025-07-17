package config

import (
	conf "github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewSetEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-environment <production|staging>",
		Short: "Configure the environment (production or staging)",
		Example: help.Examples([]string{
			"pc config set-environment production",
			"pc config set-environment staging",
		}),
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

			if secrets.OAuth2Token.Get().AccessToken != "" {
				secrets.OAuth2Token.Clear()
				msg.InfoMsg("You have been logged out; to login again, run %s", style.Code("pc login"))
			} else {
				msg.InfoMsg("To login, run %s", style.Code("pc login"))
			}

			if secrets.ApiKey.Get() != "" {
				secrets.ApiKey.Clear()
				msg.InfoMsg("API key cleared; to set a new API key, run %s", style.Code("pc config set-api-key"))
			} else {
				msg.InfoMsg("To set a new API key, run %s", style.Code("pc config set-api-key"))
			}

			if (state.TargetOrg.Get().Name != "" || state.TargetProj.Get().Name != "") && state.TargetAsst.Get().Name != "" {
				state.TargetOrg.Clear()
				state.TargetProj.Clear()
				msg.InfoMsg("Target organization and project cleared; to set a new target, run %s", style.Code("pc target -o myorg -p myproj"))
			}

			if state.TargetAsst.Get().Name != "" {
				state.TargetAsst.Clear()
				msg.InfoMsg("Target assistant cleared; to set a new target assistant, run %s", style.Code("pc assistant target -n myassistant"))
			}
		},
	}

	return cmd
}
