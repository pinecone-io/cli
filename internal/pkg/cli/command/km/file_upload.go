package km

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type UploadKnowledgeFileCmdOptions struct {
	kmName   string
	filePath string
	json     bool
}

func NewUploadKnowledgeFileCmd() *cobra.Command {
	options := UploadKnowledgeFileCmdOptions{}

	cmd := &cobra.Command{
		Use:     "file-upload",
		Short:   "Upload a file to a knowledge model",
		GroupID: help.GROUP_KM_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			targetKm := state.TargetKm.Get().Name
			if targetKm != "" {
				options.kmName = targetKm
			}
			if options.kmName == "" {
				msg.FailMsg("You must target a knowledge model or specify one with the %s flag\n", style.Emphasis("--model"))
				exit.Error(fmt.Errorf("no knowledge model specified"))
			}

			file, err := knowledge.UploadKnowledgeFile(options.kmName, options.filePath)
			if err != nil {
				msg.FailMsg("Failed to upload file %s to knowledge model %s: %s\n", style.Emphasis(options.filePath), style.Emphasis(options.kmName), err)
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(file)
				return
			}

			msg.SuccessMsg("Knowledge file %s uploaded. The file was assigned id \"%s\". \n", style.Emphasis(options.filePath), style.Emphasis(file.Id))
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.kmName, "model", "m", "", "name of the knowledge model upload a file to")
	cmd.Flags().StringVarP(&options.filePath, "file", "f", "", "the path of the file you want to upload")
	cmd.MarkFlagRequired("file")

	return cmd
}
