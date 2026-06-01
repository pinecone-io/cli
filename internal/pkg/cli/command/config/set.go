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

type SetCmdOptions struct {
	json bool
}

func NewSetCmd() *cobra.Command {
	options := SetCmdOptions{}

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
				return configKeys, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			svc := newDefaultConfigService()
			if err := runSetCmd(cmd.Context(), svc, args[0], args[1], options); err != nil {
				msg.FailJSON(options.json, "%s", err)
				exit.ErrorMsg(err.Error())
			}
		},
	}

	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}

func runSetCmd(ctx context.Context, svc ConfigService, keyName, value string, opts SetCmdOptions) error {
	// --json output for the set command
	type setOutput struct {
		Key      string   `json:"key"`
		Value    string   `json:"value"`
		Messages []string `json:"messages,omitempty"`
	}

	// Use the stored (file) value throughout so output reflects what is actually
	// persisted rather than any env var override.
	currentValue, _, err := svc.GetStored(keyName)
	if err != nil {
		return err
	}

	lines, err := svc.Set(ctx, keyName, value)
	if err != nil {
		if errors.Is(err, ErrNoChange) {
			if opts.json {
				fmt.Fprintln(os.Stdout, text.IndentJSON(setOutput{Key: keyName, Value: currentValue}))
				return nil
			}
			msg.InfoMsg("%s is already set to %s", style.Emphasis(keyName), style.Emphasis(currentValue))
			return nil
		}
		return err
	}

	if opts.json {
		storedValue, _, _ := svc.GetStored(keyName)
		fmt.Fprintln(os.Stdout, text.IndentJSON(setOutput{Key: keyName, Value: storedValue, Messages: lines}))
		return nil
	}

	msg.SuccessMsg("%s updated", style.Emphasis(keyName))
	for _, line := range lines {
		msg.InfoMsg("%s", line)
	}
	return nil
}
