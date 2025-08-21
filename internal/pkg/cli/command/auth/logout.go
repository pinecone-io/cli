package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/spf13/cobra"
)

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "Delete all saved tokens and keys",
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			secrets.ConfigFile.Clear()
			msg.SuccessMsg("API keys and user access tokens cleared.")

			state.ConfigFile.Clear()
			msg.SuccessMsg("State cleared.")
		},
	}

	return cmd
}
