package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type DescribeKnowledgeModelCmdOptions struct {
	kmName string
	json   bool
}

func NewDescribeKnowledgeModelCmd() *cobra.Command {
	options := DescribeKnowledgeModelCmdOptions{}

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a knowledge model",
		Run: func(cmd *cobra.Command, args []string) {
			model, err := knowledge.DescribeKnowledgeModel(options.kmName)
			if err != nil {
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(model)
				return
			} else {
				presenters.PrintDescribeKnowledgeModelTable(model)
			}

		},
	}
	// required flags
	cmd.Flags().StringVarP(&options.kmName, "name", "n", "", "name of the knowledge base to describe")
	cmd.MarkFlagRequired("name")

	// optional flags
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
