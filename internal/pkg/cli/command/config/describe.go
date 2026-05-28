package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeCmdOptions struct {
	reveal bool
	json   bool
}

func NewDescribeCmd() *cobra.Command {
	options := DescribeCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe <key>",
		Short: "Show detailed information about a configuration setting",
		Example: help.Examples(`
		    pc config describe api-key
		    pc config describe environment
		    pc config describe color --json
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
				type payload struct {
					Key             string   `json:"key"`
					Value           string   `json:"value"`
					Description     string   `json:"description"`
					LongDescription string   `json:"long_description,omitempty"`
					Sensitive       bool     `json:"sensitive"`
					ValidValues     []string `json:"valid_values,omitempty"`
				}
				fmt.Fprintln(os.Stdout, text.IndentJSON(payload{
					Key:             keyName,
					Value:           value,
					Description:     keyDesc.Description,
					LongDescription: keyDesc.LongDescription,
					Sensitive:       keyDesc.Sensitive,
					ValidValues:     keyDesc.ValidValues,
				}))
				return
			}

			w := presenters.NewTabWriter()
			fmt.Fprintf(w, "KEY\t%s\n", keyName)
			fmt.Fprintf(w, "VALUE\t%s\n", displayValue(value))
			fmt.Fprintf(w, "SENSITIVE\t%s\n", text.BoolToString(keyDesc.Sensitive))
			if len(keyDesc.ValidValues) > 0 {
				fmt.Fprintf(w, "VALID VALUES\t%s\n", strings.Join(keyDesc.ValidValues, ", "))
			}
			fmt.Fprintf(w, "DESCRIPTION\t%s\n", keyDesc.Description)
			w.Flush()

			if keyDesc.LongDescription != "" {
				fmt.Fprintln(os.Stdout)
				fmt.Fprintln(os.Stdout, keyDesc.LongDescription)
			}
		},
	}

	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "Reveal the full value for sensitive settings like api-key")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}
