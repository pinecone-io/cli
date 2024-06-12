package assistant

import (
	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type CreateAssistantCmdOptions struct {
	name string
	json bool
}

func NewCreateAssistantCmd() *cobra.Command {
	options := CreateAssistantCmdOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create an assistant",
		GroupID: help.GROUP_ASSISTANT_MANAGEMENT.ID,
		Run: func(cmd *cobra.Command, args []string) {
			assistant, err := assistants.CreateAssistant(options.name)
			if err != nil {
				msg.FailMsg("Failed to create assistant %s: %s\n", style.Emphasis(options.name), err)
				exit.Error(err)
			}
			msg.SuccessMsg("assistant %s created successfully.\n", style.Emphasis(assistant.Name))

			if options.json {
				text.PrettyPrintJSON(assistant)
			} else {
				presenters.PrintDescribeAssistantTable(assistant)
			}
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the assistant to create")
	cmd.MarkFlagRequired("name")
	return cmd
}
