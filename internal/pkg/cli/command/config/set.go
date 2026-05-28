package config

import (
	"errors"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Example: help.Examples(`
		    pc config set api-key pcsk_...
		    pc config set environment staging
		    pc config set color false
		`),
		Args: cobra.ExactArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return configKeyOrder, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			keyName, value := args[0], args[1]

			keyDesc, err := lookupKey(keyName)
			if err != nil {
				msg.FailMsg("%s", err)
				exit.ErrorMsg(err.Error())
				return
			}

			oldVal := keyDesc.getStr()

			if err := keyDesc.setStr(value); err != nil {
				if errors.Is(err, ErrNoChange) {
					msg.InfoMsg("%s is already set to %s", style.Emphasis(keyName), style.Emphasis(oldVal))
					return
				}
				msg.FailMsg("%s", err)
				exit.ErrorMsg(err.Error())
				return
			}

			msg.SuccessMsg("%s updated", style.Emphasis(keyName))

			if keyDesc.onChange != nil {
				for _, line := range keyDesc.onChange(cmd.Context(), oldVal, value) {
					msg.InfoMsg("%s", line)
				}
			}
		},
	}

	return cmd
}
