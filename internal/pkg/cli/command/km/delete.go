package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/spf13/cobra"
)

type DeleteKnowledgeModelOptions struct {
	kbName string
	json   bool
}

func NewDeleteKnowledgeModelCmd() *cobra.Command {
	options := DeleteKnowledgeModelOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a knowledge model",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := knowledge.DeleteKnowledgeModel(options.kbName)
			if err != nil {
				exit.Error(err)
			}

			msg.SuccessMsg("Knowledge model %s deleted.\n", options.kbName)
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.kbName, "name", "n", "", "name of the knowledge model to delete")
	cmd.MarkFlagRequired("name")
	return cmd
}
