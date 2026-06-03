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
		ValidArgs: visibleKeys(),
		Run: func(cmd *cobra.Command, args []string) {
			svc := newDefaultConfigService()
			if err := runDescribeCmd(svc, args[0], options); err != nil {
				msg.FailJSON(options.json, "%s", err)
				exit.ErrorMsg(err.Error())
			}
		},
	}

	cmd.Flags().BoolVar(&options.reveal, "reveal", false, "Reveal the full value for sensitive settings like api-key")
	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}

func runDescribeCmd(svc ConfigService, keyName string, opts DescribeCmdOptions) error {
	// --json output for the describe command
	type describeOutput struct {
		Key            string   `json:"key"`
		Value          string   `json:"value"`
		EnvVarName     string   `json:"env_var_name,omitempty"`
		EnvVarOverride *bool    `json:"env_var_override,omitempty"`
		Description    string   `json:"description"`
		Sensitive      bool     `json:"sensitive"`
		ValidValues    []string `json:"valid_values,omitempty"`
	}

	desc, err := svc.Describe(keyName)
	if err != nil {
		return err
	}

	value := desc.Value
	if desc.Sensitive && !opts.reveal {
		value = presenters.MaskHeadTail(value, 4, 4)
	}
	if opts.json {
		out := describeOutput{
			Key:         desc.Key,
			Value:       value,
			EnvVarName:  desc.EnvVarName,
			Description: desc.Description,
			Sensitive:   desc.Sensitive,
			ValidValues: desc.ValidValues,
		}
		if desc.EnvVarName != "" {
			out.EnvVarOverride = &desc.EnvVarOverride
		}
		fmt.Fprintln(os.Stdout, text.IndentJSON(out))
		return nil
	}

	w := presenters.NewTabWriter()
	fmt.Fprintf(w, "KEY\t%s\n", desc.Key)
	fmt.Fprintf(w, "VALUE\t%s\n", displayValue(value))
	if desc.EnvVarName != "" {
		fmt.Fprintf(w, "ENV VAR NAME\t$%s\n", desc.EnvVarName)
		fmt.Fprintf(w, "ENV VAR OVERRIDE\t%s\n", text.BoolToString(desc.EnvVarOverride))
	}
	fmt.Fprintf(w, "SENSITIVE\t%s\n", text.BoolToString(desc.Sensitive))
	if len(desc.ValidValues) > 0 {
		fmt.Fprintf(w, "VALID VALUES\t%s\n", strings.Join(desc.ValidValues, ", "))
	}
	fmt.Fprintf(w, "DESCRIPTION\t%s\n", desc.Description)
	w.Flush()

	if desc.LongDescription != "" {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, desc.LongDescription)
	}

	return nil
}
