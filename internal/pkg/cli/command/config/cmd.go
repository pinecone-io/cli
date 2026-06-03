package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/spf13/cobra"
)

var (
	configHelp = help.LongF(`
		Manage configuration for the Pinecone CLI.

		Configuration for this CLI is stored in a file called config.yaml in the %s directory.
	`, configuration.NewConfigLocations().ConfigPath)
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration for the Pinecone CLI",
		Long:  configHelp,
	}

	// Primary commands
	cmd.AddCommand(NewGetCmd())
	cmd.AddCommand(NewSetCmd())
	cmd.AddCommand(NewUnsetCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewDescribeCmd())

	// Deprecated aliases kept for backwards compatibility
	cmd.AddCommand(NewGetApiKeyCmd())
	cmd.AddCommand(NewSetApiKeyCmd())
	cmd.AddCommand(NewSetColorCmd())
	cmd.AddCommand(NewSetEnvCmd())

	return cmd
}
