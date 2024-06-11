package assistant

import (
	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type DeleteKnowledgeFileCmdOptions struct {
	name   string
	fileId string
}

func NewDeleteKnowledgeFileCmd() *cobra.Command {
	options := DeleteKnowledgeFileCmdOptions{}

	cmd := &cobra.Command{
		Use:     "file-delete",
		Short:   "Delete a file in an assistant",
		GroupID: help.GROUP_ASSISTANT_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			targetKm := state.TargetAsst.Get().Name
			if targetKm != "" {
				options.name = targetKm
			}
			if options.name == "" {
				msg.FailMsg("You must target an assistant or specify one with the %s flag\n", style.Emphasis("--model"))
				exit.ErrorMsg("no assistant specified")
			}

			_, err := assistants.DeleteKnowledgeFile(options.name, options.fileId)
			if err != nil {
				msg.FailMsg("Failed to delete file %s in assistant %s: %s\n", style.Emphasis(options.fileId), style.Emphasis(options.name), err)
				exit.Error(err)
			}

			msg.SuccessMsg("Assistant file %s deleted.\n", style.Emphasis(options.fileId))
		},
	}

	cmd.Flags().StringVarP(&options.name, "model", "m", "", "name of the assistant to list files for")
	cmd.Flags().StringVarP(&options.fileId, "id", "i", "", "id of the file to describe")
	cmd.MarkFlagRequired("id")

	return cmd
}
