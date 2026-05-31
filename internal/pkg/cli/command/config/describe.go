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
				msg.FailMsg("%s", err)
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
		Key             string   `json:"key"`
		Value           string   `json:"value"`
		Description     string   `json:"description"`
		LongDescription string   `json:"long_description,omitempty"`
		Sensitive       bool     `json:"sensitive"`
		ValidValues     []string `json:"valid_values,omitempty"`
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
		fmt.Fprintln(os.Stdout, text.IndentJSON(describeOutput{
			Key:             desc.Key,
			Value:           value,
			Description:     desc.Description,
			LongDescription: desc.LongDescription,
			Sensitive:       desc.Sensitive,
			ValidValues:     desc.ValidValues,
		}))
		return nil
	}

	w := presenters.NewTabWriter()
	fmt.Fprintf(w, "KEY\t%s\n", desc.Key)
	fmt.Fprintf(w, "VALUE\t%s\n", displayValue(value))
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
