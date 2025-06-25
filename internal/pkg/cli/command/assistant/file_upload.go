package assistant

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type UploadAssistantCmdOptions struct {
	assistant string
	filePath  string
	json      bool
}

func NewUploadAssistantFileCmd() *cobra.Command {
	options := UploadAssistantCmdOptions{}

	cmd := &cobra.Command{
		Use:     "file-upload",
		Short:   "Upload a file to an assistant",
		GroupID: help.GROUP_ASSISTANT_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			targetAsst := state.TargetAsst.Get().Name
			if targetAsst != "" {
				options.assistant = targetAsst
			}
			if options.assistant == "" {
				msg.FailMsg("You must target an assistant or specify one with the %s flag\n", style.Emphasis("--assistant"))
				exit.Error(fmt.Errorf("no assistant specified"))
			}

			file, err := assistants.UploadAssistantFile(options.assistant, options.filePath)
			if err != nil {
				msg.FailMsg("Failed to upload file %s to assistant %s: %s\n", style.Emphasis(options.filePath), style.Emphasis(options.assistant), err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(file)
				pcio.Println(json)
				return
			}

			msg.SuccessMsg("Assistant file %s uploaded. The file was assigned id \"%s\". \n", style.Emphasis(options.filePath), style.Emphasis(file.Id))
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.assistant, "assistant", "a", "", "name of the assistant to upload a file to")
	cmd.Flags().StringVarP(&options.filePath, "file", "f", "", "the path of the file you want to upload")
	cmd.MarkFlagRequired("file")

	return cmd
}
