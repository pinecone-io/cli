package config

import (
	"fmt"
	"os"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
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
			svc := newDefaultConfigService()
			if err := runListCmd(svc, options); err != nil {
				msg.FailJSON(options.json, "%s", err)
				exit.ErrorMsg(err.Error())
			}
		},
	}

	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "Reveal the full value for sensitive settings like api-key")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}

func runListCmd(svc ConfigService, opts ListCmdOptions) error {
	// --json output for the list command
	type listOutput struct {
		Key         string `json:"key"`
		Value       string `json:"value"`
		Description string `json:"description"`
	}

	entries := svc.List()

	if opts.json {
		jsonEntries := make([]listOutput, 0, len(entries))
		for _, e := range entries {
			value := e.Value
			if e.Sensitive && !opts.reveal {
				value = presenters.MaskHeadTail(value, 4, 4)
			}
			jsonEntries = append(jsonEntries, listOutput{Key: e.Key, Value: value, Description: e.Description})
		}
		fmt.Fprintln(os.Stdout, text.IndentJSON(jsonEntries))
		return nil
	}

	w := presenters.NewTabWriter()
	fmt.Fprintln(w, "KEY\tVALUE\tDESCRIPTION")
	for _, e := range entries {
		value := e.Value
		if e.Sensitive && !opts.reveal {
			value = presenters.MaskHeadTail(value, 4, 4)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", e.Key, displayValue(value), e.Description)
	}
	w.Flush()
	return nil
}
