package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewUnsetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unset <key>",
		Short: "Reset a configuration value to its default",
		Example: help.Examples(`
		    pc config unset api-key
		    pc config unset color
		`),
		Args:      cobra.ExactArgs(1),
		ValidArgs: configKeyOrder,
		Run: func(cmd *cobra.Command, args []string) {
			keyName := args[0]
			keyDesc, err := lookupKey(keyName)
			if err != nil {
				msg.FailMsg("%s", err)
				exit.ErrorMsg(err.Error())
				return
			}

			keyDesc.clearStr()
			msg.SuccessMsg("%s cleared", style.Emphasis(keyName))
		},
	}

	return cmd
}
