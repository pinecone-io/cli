package config

import (
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type ListCmdOptions struct {
	reveal bool
	json   bool
}

func NewListCmd() *cobra.Command {
	options := ListCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration settings and their current values",
		Example: help.Examples(`
		    pc config list
		    pc config list --reveal
		    pc config list --json
		`),
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if options.json {
				type entry struct {
					Key         string `json:"key"`
					Value       string `json:"value"`
					Description string `json:"description"`
				}
				entries := make([]entry, 0, len(configKeyOrder))
				for _, keyName := range configKeyOrder {
					keyDesc := configRegistry[keyName]
					value := keyDesc.getStr()
					if keyDesc.Sensitive && !options.reveal {
						value = presenters.MaskHeadTail(value, 4, 4)
					}
					entries = append(entries, entry{
						Key:         keyName,
						Value:       value,
						Description: keyDesc.Description,
					})
				}
				fmt.Fprintln(os.Stdout, text.IndentJSON(entries))
				return
			}

			w := presenters.NewTabWriter()
			fmt.Fprintln(w, "KEY\tVALUE\tDESCRIPTION")
			for _, keyName := range configKeyOrder {
				keyDesc := configRegistry[keyName]
				value := keyDesc.getStr()
				if keyDesc.Sensitive && !options.reveal {
					value = presenters.MaskHeadTail(value, 4, 4)
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", keyName, displayValue(value), keyDesc.Description)
			}
			w.Flush()
		},
	}

	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "Reveal the full value for sensitive settings like api-key")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}
