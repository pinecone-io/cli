package config

import (
	conf "github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

func NewSetStagingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-staging",
		Short: "Configure whether the CLI uses staging environments",
		Run: func(cmd *cobra.Command, args []string) {
			stagingArg := args[0]

			var stagingSetting bool
			switch stagingArg {
			case "true", "on", "1":
				stagingSetting = true
			default:
				stagingSetting = false
			}

			conf.Staging.Set(stagingSetting)
			pcio.Printf("Config property %s updated to %s\n", style.Emphasis("staging"), style.Emphasis(text.BoolToString(stagingSetting)))
		},
	}

	return cmd
}
