package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/spf13/cobra"
)

type DeleteKnowledgeFileCmdOptions struct {
	kmName string
	fileId string
	json   bool
}

func NewDeleteKnowledgeFileCmd() *cobra.Command {
	options := DeleteKnowledgeFileCmdOptions{}

	cmd := &cobra.Command{
		Use:   "delete-file",
		Short: "Delete a file in a knowledge model",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := knowledge.DeleteKnowledgeFile(options.kmName, options.fileId)
			if err != nil {
				exit.Error(err)
			}

			msg.SuccessMsg("Knowledge file %s deleted.\n", options.fileId)
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.kmName, "name", "n", "", "name of the knowledge model to list files for")
	cmd.Flags().StringVarP(&options.fileId, "id", "i", "", "id of the file to describe")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("id")

	return cmd
}
