package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewSetApiKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-api-key",
		Short: "Manually set the global API key for the Pinecone CLI",
		Example: help.Examples(`
		    pc config set-api-key <api-key>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			newApiKey := args[0]
			secrets.GlobalApiKey.Set(newApiKey)
			msg.SuccessMsg("Config property %s updated.", style.Emphasis("api_key"))
			msg.InfoMsg("To clear the global API key, run %s.", style.Code("pc auth clear --global-api-key"))
		},
	}

	return cmd
}
