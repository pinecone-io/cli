package auth

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/auth"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/spf13/cobra"
)

func NewCmdAuthStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Show the current authentication status of the Pinecone CLI",
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runAuthStatus(cmd); err != nil {
				log.Error().Err(err).Msg("Error retrieving authentication status")
				exit.Error(pcio.Errorf("error retrieving authentication status: %w", err))
			}
		},
	}
	return cmd
}

func runAuthStatus(cmd *cobra.Command) error {
	token, err := auth.Token(cmd.Context())
	if err != nil { // This should only error on a network request to refresh the token
		return err
	}
	presenters.PrintAuthStatus(token)

	return nil
}
