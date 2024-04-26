package config

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/spf13/cobra"
)

var configHelpText = fmt.Sprintf(`Configuration for this CLI is stored in a file called 
config.yaml in the %s directory.`, configuration.NewConfigLocations().ConfigPath)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage configuration for the Pinecone CLI",
		Long:  configHelpText,
	}

	cmd.AddCommand(NewSetColorCmd())
	cmd.AddCommand(NewSetApiKeyCmd())

	return cmd
}
