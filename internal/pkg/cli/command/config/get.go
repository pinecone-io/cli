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
		ValidArgs: configKeyOrder,
		Run: func(cmd *cobra.Command, args []string) {
			keyName := args[0]
			keyDesc, err := lookupKey(keyName)
			if err != nil {
				msg.FailMsg("%s", err)
				exit.ErrorMsg(err.Error())
				return
			}

			value := keyDesc.getStr()
			if keyDesc.Sensitive && !options.reveal {
				value = presenters.MaskHeadTail(value, 4, 4)
			}

			if options.json {
				fmt.Fprintln(os.Stdout, text.IndentJSON(struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				}{Key: keyName, Value: value}))
				return
			}

			msg.InfoMsg("%s: %s", style.Emphasis(keyName), displayValue(value))
		},
	}

	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "Reveal the full value for sensitive settings like api-key")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}
