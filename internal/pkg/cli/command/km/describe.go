package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
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
