package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
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
				pcio.Printf("You must target a knowledge model or specify one with the %s flag\n", style.Emphasis("--name"))
				return
			}

			// Check if file is pdf or txt
			if !knowledge.IsSupportedFile(options.filePath) {
				pcio.Printf("File type not supported. Supported file types are .pdf and .txt\n")
				return
			}

			file, err := knowledge.UploadKnowledgeFile(options.kmName, options.filePath)
			if err != nil {
				exit.Error(err)
			}

			msg.SuccessMsg("Knowledge file %s uploaded. ID=%s \n", options.filePath, file.Id)
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.kmName, "model", "m", "", "name of the knowledge model upload a file to")
	cmd.Flags().StringVarP(&options.filePath, "file", "f", "", "the path of the file you want to upload")
	cmd.MarkFlagRequired("file")

	return cmd
}
