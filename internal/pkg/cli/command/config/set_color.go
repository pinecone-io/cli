package config

import (
	conf "github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

func NewSetColorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-color",
		Short: "Configure whether the CLI prints output with color",
		Example: help.Examples([]string{
			"pc config set-color true",
			"pc config set-color false",
		}),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				msg.FailMsg("Please provide a value for color. Accepted values are 'true', 'false'")
				exit.ErrorMsg("No value provided for color")
			}

			colorArg := args[0]

			var colorSetting bool
			switch colorArg {
			case "true", "on", "1":
				colorSetting = true
			default:
				colorSetting = false
			}

			conf.Color.Set(colorSetting)
			msg.SuccessMsg("Config property %s updated to %s\n", style.Emphasis("color"), style.Emphasis(text.BoolToString(colorSetting)))
		},
	}

	return cmd
}
