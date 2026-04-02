package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/login"
	"github.com/spf13/cobra"
)

// NewDaemonCmd returns the hidden internal command spawned by `pc login --json`
// to run the OAuth callback server in the background.
func NewDaemonCmd() *cobra.Command {
	var sessionId string

	cmd := &cobra.Command{
		Use:    "_daemon",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			login.RunDaemon(sessionId)
		},
	}

	cmd.Flags().StringVar(&sessionId, "session-id", "", "session ID for the pending auth flow")
	_ = cmd.MarkFlagRequired("session-id")

	return cmd
}
