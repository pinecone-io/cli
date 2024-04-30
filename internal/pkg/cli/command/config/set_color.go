package config

import (
	conf "github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

func NewSetColorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-color",
		Short: "Configure whether the CLI prints output with color",
		Run: func(cmd *cobra.Command, args []string) {
			colorArg := args[0]

			var colorSetting bool
			switch colorArg {
			case "true", "on", "1":
				colorSetting = true
			default:
				colorSetting = false
			}

			conf.Color.Set(colorSetting)
			pcio.Printf("Config property %s updated to %s\n", style.Emphasis("color"), style.Emphasis(text.BoolToString(colorSetting)))
		},
	}

	return cmd
}
