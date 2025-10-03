package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage configuration for the Pinecone CLI",
		Long: help.LongF(`Configuration for this CLI is stored in a file called config.yaml in the %s directory.`,
			configuration.NewConfigLocations().ConfigPath),
	}

	cmd.AddCommand(NewSetColorCmd())
	cmd.AddCommand(NewSetApiKeyCmd())
	cmd.AddCommand(NewGetApiKeyCmd())
	cmd.AddCommand(NewSetEnvCmd())

	return cmd
}
