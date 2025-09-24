package auth

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/spf13/cobra"
)

type ClearCmdOptions struct {
	serviceAccount bool
	globalAPIKey   bool
}

func NewClearCmd() *cobra.Command {
	options := ClearCmdOptions{}

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Allows you to clear a configured service account (client id and secret), or global API key",
		Example: heredoc.Doc(`
		# Clear configured service account credentials
		$ pc auth clear --service-account

		# Clear configured global API key
		$ pc auth clear --global-api-key

		# Clear both configured service account credentials and global API key
		$ pc auth clear --service-account --global-api-key
		`),
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			if !options.serviceAccount && !options.globalAPIKey {
				msg.FailMsg("Please specify either --service-account or --global-api-key")
				exit.ErrorMsg("No option specified")
			}

			if options.serviceAccount {
				secrets.ClientId.Clear()
				secrets.ClientSecret.Clear()
				msg.SuccessMsg("Service account (client id and secret) cleared")
			}

			if options.globalAPIKey {
				secrets.GlobalApiKey.Clear()
				msg.SuccessMsg("Global API key cleared")
			}
		},
	}

	cmd.Flags().BoolVar(&options.serviceAccount, "service-account", false, "Clear the configured service account (client id and secret)")
	cmd.Flags().BoolVar(&options.globalAPIKey, "global-api-key", false, "Clear the configured global API key")

	return cmd
}
