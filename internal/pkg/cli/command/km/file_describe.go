package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeKnowledgeFileOptions struct {
	kmName string
	fileId string
	json   bool
}

func NewDescribeKnowledgeFileCmd() *cobra.Command {
	options := DescribeKnowledgeFileOptions{}

	cmd := &cobra.Command{
		Use:   "describe-file",
		Short: "Describe a file in a knowledge model",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := knowledge.DescribeKnowledgeModelFile(options.kmName, options.fileId)
			if err != nil {
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(file)
			} else {
				presenters.PrintDescribeKnowledgeFileTable(file)
			}
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.kmName, "name", "n", "", "name of the knowledge model to list files for")
	cmd.Flags().StringVarP(&options.fileId, "id", "i", "", "id of the file to describe")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("id")

	return cmd
}
