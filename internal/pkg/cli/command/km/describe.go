package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
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
		Use:     "describe",
		Short:   "Describe a knowledge model",
		GroupID: help.GROUP_KM_MANAGEMENT.ID,
		Run: func(cmd *cobra.Command, args []string) {
			// If no name is provided, use the target knowledge model
			if options.kmName == "" {
				targetKm := state.TargetKm.Get().Name
				options.kmName = targetKm
			}
			if options.kmName == "" {
				msg.FailMsg("You must target a knowledge model or specify one with the %s flag\n", style.Emphasis("--name"))
				exit.ErrorMsg("No knowledge model specified")
				return
			}

			model, err := knowledge.DescribeKnowledgeModel(options.kmName)
			if err != nil {
				msg.FailMsg("Failed to describe knowledge model %s: %s\n", style.Emphasis(options.kmName), err)
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

	cmd.Flags().StringVarP(&options.kmName, "name", "n", "", "name of the knowledge base to describe")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}
