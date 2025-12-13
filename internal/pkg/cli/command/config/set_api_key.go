package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

var (
	setAPIKeyHelp = help.Long(`
		Configure the CLI to authenticate with Pinecone using an API key.

		This overrides any target context set through user login or service account credentials.
		To clear the explicit API key, run 'pc auth clear --api-key'.
	`)
)

func NewSetApiKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-api-key",
		Short: "Configure the CLI to authenticate with Pinecone using a default API key",
		Long:  setAPIKeyHelp,
		Example: help.Examples(`
		    pc config set-api-key "api-key-value"
		`),
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				msg.FailMsg("Please provide an API key value.")
				exit.ErrorMsg("No value provided for API key")
			}

			newApiKey := args[0]
			if newApiKey == "" {
				msg.FailMsg("Please provide an API key value.")
				exit.ErrorMsg("No value provided for API key")
			}

			secrets.DefaultAPIKey.Set(newApiKey)
			msg.SuccessMsg("Config property %s updated.", style.Emphasis("api_key"))
			msg.InfoMsg("To clear the default API key, run %s.", style.Code("pc auth clear --api-key"))
		},
	}

	return cmd
}
