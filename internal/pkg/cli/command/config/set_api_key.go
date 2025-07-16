package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewSetApiKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-api-key",
		Short: "Manually set the API key for the Pinecone CLI",
		Run: func(cmd *cobra.Command, args []string) {
			newApiKey := args[0]
			secrets.ApiKey.Set(newApiKey)
			msg.SuccessMsg("Config property %s updated.", style.Emphasis("api_key"))
			msg.InfoMsg("To clear saved keys, run %s.", style.Code("pc logout"))
		},
	}

	return cmd
}
