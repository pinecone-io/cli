package config

import (
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type GetCmdOptions struct {
	reveal bool
	json   bool
}

func NewGetCmd() *cobra.Command {
	options := GetCmdOptions{}

	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get the current value of a configuration setting",
		Example: help.Examples(`
		    pc config get api-key
		    pc config get api-key --reveal
		    pc config get environment
		    pc config get color
		`),
		Args:      cobra.ExactArgs(1),
		ValidArgs: visibleKeys(),
		Run: func(cmd *cobra.Command, args []string) {
			svc := newDefaultConfigService()
			if err := runGetCmd(svc, args[0], options); err != nil {
				msg.FailMsg("%s", err)
				exit.ErrorMsg(err.Error())
			}
		},
	}

	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "Reveal the full value for sensitive settings like api-key")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}

func runGetCmd(svc ConfigService, keyName string, opts GetCmdOptions) error {
	// --json output for the get command
	type getOutput struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	value, sensitive, err := svc.Get(keyName)
	if err != nil {
		return err
	}

	if sensitive && !opts.reveal {
		value = presenters.MaskHeadTail(value, 4, 4)
	}

	if opts.json {
		fmt.Fprintln(os.Stdout, text.IndentJSON(getOutput{Key: keyName, Value: value}))
		return nil
	}

	msg.InfoMsg("%s: %s", style.Emphasis(keyName), displayValue(value))
	return nil
}
