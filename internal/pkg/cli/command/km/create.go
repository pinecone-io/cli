package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

type CreateKnowledgeModelCmdOptions struct {
	name string
	json bool
}

func NewCreateKnowledgeModelCmd() *cobra.Command {
	options := CreateKnowledgeModelCmdOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a knowledge model",
		Run: func(cmd *cobra.Command, args []string) {
			model, err := knowledge.CreateKnowledgeModel(options.name)
			if err != nil {
				msg.FailMsg("Failed to create knowledge model %s: %s\n", style.Emphasis(options.name), err)
				exit.Error(err)
			}
			msg.SuccessMsg("Knowledge model %s created successfully.\n", style.Emphasis(model.Name))
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the knowledge model")
	cmd.MarkFlagRequired("name")
	return cmd
}
