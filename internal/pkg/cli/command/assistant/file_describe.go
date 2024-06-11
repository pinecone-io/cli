package assistant

import (
	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeKnowledgeFileCmdOptions struct {
	name   string
	fileId string
	json   bool
}

func NewDescribeKnowledgeFileCmd() *cobra.Command {
	options := DescribeKnowledgeFileCmdOptions{}

	cmd := &cobra.Command{
		Use:     "file-describe",
		Short:   "Describe a file in an assistant",
		GroupID: help.GROUP_ASSISTANT_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			targetKm := state.TargetAsst.Get().Name
			if targetKm != "" {
				options.name = targetKm
			}
			if options.name == "" {
				pcio.Printf("You must target an assistant or specify one with the %s flag\n", style.Emphasis("--model"))
				return
			}

			file, err := assistants.DescribeAssistantFile(options.name, options.fileId)
			if err != nil {
				msg.FailMsg("Failed to describe file %s in assistant: %s\n", style.Emphasis(options.fileId), err)
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(file)
			} else {
				presenters.PrintDescribeAssistantFileTable(file)
			}
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.name, "model", "m", "", "name of the assistant to list files for")
	cmd.Flags().StringVarP(&options.fileId, "id", "i", "", "id of the file to describe")
	cmd.MarkFlagRequired("id")

	return cmd
}