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
	all    bool
}

func NewListCmd() *cobra.Command {
	options := ListCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration settings and their current values",
		Example: help.Examples(`
		    pc config list
		    pc config list --all
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
	cmd.Flags().BoolVarP(&options.all, "all", "a", false, "Include hidden settings such as environment")

	return cmd
}

func runListCmd(svc ConfigService, opts ListCmdOptions) error {
	// --json output for the list command
	type listOutput struct {
		Key            string `json:"key"`
		Value          string `json:"value"`
		EnvVarName     string `json:"env_var_name,omitempty"`
		EnvVarOverride *bool  `json:"env_var_override,omitempty"`
		Description    string `json:"description"`
		Hidden         bool   `json:"hidden,omitempty"`
	}

	entries := svc.List(opts.all)

	if opts.json {
		jsonEntries := make([]listOutput, 0, len(entries))
		for _, e := range entries {
			value := e.Value
			if e.Sensitive && !opts.reveal {
				value = presenters.MaskHeadTail(value, 4, 4)
			}
			entry := listOutput{Key: e.Key, Value: value, EnvVarName: e.EnvVarName, Description: e.Description, Hidden: e.Hidden}
			if e.EnvVarName != "" {
				entry.EnvVarOverride = &e.EnvVarOverride
			}
			jsonEntries = append(jsonEntries, entry)
		}
		fmt.Fprintln(os.Stdout, text.IndentJSON(jsonEntries))
		return nil
	}

	w := presenters.NewTabWriter()
	fmt.Fprintln(w, "KEY\tVALUE\tENV VAR NAME\tENV VAR OVERRIDE\tDESCRIPTION")
	for _, e := range entries {
		value := e.Value
		if e.Sensitive && !opts.reveal {
			value = presenters.MaskHeadTail(value, 4, 4)
		}
		displayVal := displayValue(value)
		if e.EnvVarOverride {
			displayVal = fmt.Sprintf("%s [$%s]", displayVal, e.EnvVarName)
		}

		envOverride := ""
		if e.EnvVarName != "" {
			envOverride = text.BoolToString(e.EnvVarOverride)
		}

		fmt.Fprintf(w,
			"%s\t%s\t%s\t%s\t%s\n",
			e.Key, displayVal,
			e.EnvVarName,
			envOverride,
			e.Description)
	}
	w.Flush()
	return nil
}
