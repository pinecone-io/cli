package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type DeleteKnowledgeModelCmdOptions struct {
	kmName string
	json   bool
}

func NewDeleteKnowledgeModelCmd() *cobra.Command {
	options := DeleteKnowledgeModelCmdOptions{}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a knowledge model",
		GroupID: help.GROUP_KM_MANAGEMENT.ID,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := knowledge.DeleteKnowledgeModel(options.kmName)
			if err != nil {
				msg.FailMsg("Failed to delete knowledge model %s: %s\n", style.Emphasis(options.kmName), err)
				exit.Error(err)
			}

			msg.SuccessMsg("Knowledge model %s deleted.\n", style.Emphasis(options.kmName))
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.kmName, "name", "n", "", "name of the knowledge model to delete")
	cmd.MarkFlagRequired("name")
	return cmd
}
