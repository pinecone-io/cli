package assistant

import (
	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeAssistantCmdOptions struct {
	name string
	json bool
}

func NewDescribeAssistantCmd() *cobra.Command {
	options := DescribeAssistantCmdOptions{}

	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Describe an assistant",
		GroupID: help.GROUP_ASSISTANT_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			// If no name is provided, use the target assistant
			if options.name == "" {
				targetAsst := state.TargetAsst.Get().Name
				options.name = targetAsst
			}
			if options.name == "" {
				msg.FailMsg("You must target an assistant or specify one to describe with the %s flag\n", style.Emphasis("--name"))
				exit.ErrorMsg("No assistant specified")
				return
			}

			assistant, err := assistants.DescribeAssistant(options.name)
			if err != nil {
				msg.FailMsg("Failed to describe assistant %s: %s\n", style.Emphasis(options.name), err)
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(assistant)
				return
			} else {
				presenters.PrintDescribeAssistantTable(assistant)
			}

		},
	}

	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the assistant to describe")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
