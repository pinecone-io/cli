package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type UnsetCmdOptions struct {
	json bool
}

func NewUnsetCmd() *cobra.Command {
	options := UnsetCmdOptions{}

	cmd := &cobra.Command{
		Use:   "unset <key>",
		Short: "Reset a configuration value to its default",
		Example: help.Examples(`
		    pc config unset api-key
		    pc config unset color
		`),
		Args:      cobra.ExactArgs(1),
		ValidArgs: visibleKeys(),
		Run: func(cmd *cobra.Command, args []string) {
			svc := newDefaultConfigService()
			if err := runUnsetCmd(cmd.Context(), svc, args[0], options); err != nil {
				msg.FailJSON(options.json, "%s", err)
				exit.ErrorMsg(err.Error())
			}
		},
	}

	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}

func runUnsetCmd(ctx context.Context, svc ConfigService, keyName string, opts UnsetCmdOptions) error {
	// --json output for the unset command
	type unsetOutput struct {
		Key     string `json:"key"`
		Cleared bool   `json:"cleared"`
	}

	lines, err := svc.Unset(ctx, keyName)
	if err != nil {
		if errors.Is(err, ErrNoChange) {
			if opts.json {
				fmt.Fprintln(os.Stdout, text.IndentJSON(unsetOutput{Key: keyName, Cleared: false}))
				return nil
			}
			msg.InfoMsg("%s is already at its default value", style.Emphasis(keyName))
			return nil
		}
		return err
	}

	if opts.json {
		fmt.Fprintln(os.Stdout, text.IndentJSON(unsetOutput{Key: keyName, Cleared: true}))
		return nil
	}

	msg.SuccessMsg("%s cleared", style.Emphasis(keyName))
	for _, line := range lines {
		msg.InfoMsg("%s", line)
	}
	return nil
}
