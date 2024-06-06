package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
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
		Use:   "upload-file",
		Short: "Upload a file to a knowledge model",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := knowledge.UploadKnowledgeFile(options.kmName, options.filePath)
			if err != nil {
				exit.Error(err)
			}

			msg.SuccessMsg("Knowledge file %s uploaded. ID=%s \n", options.filePath, file.Id)
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.kmName, "name", "n", "", "name of the knowledge model upload a file to")
	cmd.Flags().StringVarP(&options.filePath, "file", "f", "", "the path of the file you want to upload")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("file")

	return cmd
}
