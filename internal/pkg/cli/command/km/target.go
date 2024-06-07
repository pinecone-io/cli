package km

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type KmTargetCmdOptions struct {
	KmName string
	json   bool
	clear  bool
}

var kmTargetHelpPart1 string = text.WordWrap(`There are many knowledge model commands which target a specific
knowledge model. This command allows you to set and clear the target knowledge model for performing operations.`, 80)

var targetHelp = pcio.Sprintf("%s\n", kmTargetHelpPart1)

func NewKmTargetCmd() *cobra.Command {
	options := KmTargetCmdOptions{}

	cmd := &cobra.Command{
		Use:     "target <flags>",
		Short:   "Set the target knowledge model",
		Long:    targetHelp,
		GroupID: help.GROUP_KM_TARGETING.ID,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug().
				Str("kmName", options.KmName).
				Bool("json", options.json).
				Bool("clear", options.clear).
				Msg("km target command invoked")

			// Clear targets
			if options.clear {
				state.ConfigFile.Clear()
				pcio.Print("target knowledge model cleared")
				return
			}

			// Print current target if no knowledge model is specified
			if options.KmName == "" {
				pcio.Printf("Current target knowledge model: %s\n", state.TargetKm.Get().Name)

				pcio.Printf("To target a knowledge model, use %s \n\n", style.Code("pinecone km target --name <name>"))
				presenters.PrintTargetKnowledgeModel(state.GetTargetContext())
				return
			}

			// Set target knowledge model
			if options.KmName != "" {
				state.TargetKm.Set(&state.TargetKnowledgeModel{Name: options.KmName})
			}

			pcio.Print("Target knowledge model set")
		},
	}

	cmd.Flags().StringVarP(&options.KmName, "name", "n", "", "name of the knowledge model to target")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().BoolVar(&options.clear, "clear", false, "clear the target knowledge model")

	return cmd
}
