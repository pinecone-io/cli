package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

var configHelpText = text.WordWrap(pcio.Sprintf(`Configuration for this CLI is stored in a file called 
config.yaml in the %s directory.`, configuration.NewConfigLocations().ConfigPath), 80)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage configuration for the Pinecone CLI",
		Long:  configHelpText,
	}

	// TODO - Remove NewSetStagingCmd(). Look at adding more robust support for configuration through this command.

	cmd.AddCommand(NewSetColorCmd())
	cmd.AddCommand(NewSetApiKeyCmd())
	cmd.AddCommand(NewGetApiKeyCmd())
	cmd.AddCommand(NewSetStagingCmd())

	return cmd
}
