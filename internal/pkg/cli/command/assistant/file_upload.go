package assistant

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type UploadAssistantCmdOptions struct {
	name     string
	filePath string
	json     bool
}

func NewUploadAssistantFileCmd() *cobra.Command {
	options := UploadAssistantCmdOptions{}

	cmd := &cobra.Command{
		Use:     "file-upload",
		Short:   "Upload a file to an assistant",
		GroupID: help.GROUP_ASSISTANT_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			targetKm := state.TargetAsst.Get().Name
			if targetKm != "" {
				options.name = targetKm
			}
			if options.name == "" {
				msg.FailMsg("You must target an assistant or specify one with the %s flag\n", style.Emphasis("--assistant"))
				exit.Error(fmt.Errorf("no assistant specified"))
			}

			file, err := assistants.UploadAssistantFile(options.name, options.filePath)
			if err != nil {
				msg.FailMsg("Failed to upload file %s to assistant %s: %s\n", style.Emphasis(options.filePath), style.Emphasis(options.name), err)
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(file)
				return
			}

			msg.SuccessMsg("Assistant file %s uploaded. The file was assigned id \"%s\". \n", style.Emphasis(options.filePath), style.Emphasis(file.Id))
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.name, "assistant", "a", "", "name of the assistant to upload a file to")
	cmd.Flags().StringVarP(&options.filePath, "file", "f", "", "the path of the file you want to upload")
	cmd.MarkFlagRequired("file")

	return cmd
}
