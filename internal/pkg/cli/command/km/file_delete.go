package km

import (
	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
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
		Use:     "file-delete",
		Short:   "Delete a file in a knowledge model",
		GroupID: help.GROUP_KM_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			targetKm := state.TargetKm.Get().Name
			if targetKm != "" {
				options.kmName = targetKm
			}
			if options.kmName == "" {
				pcio.Printf("You must target a knowledge model or specify one with the %s flag\n", style.Emphasis("--model"))
				return
			}

			_, err := knowledge.DeleteKnowledgeFile(options.kmName, options.fileId)
			if err != nil {
				exit.Error(err)
			}

			msg.SuccessMsg("Knowledge file %s deleted.\n", options.fileId)
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.kmName, "model", "m", "", "name of the knowledge model to list files for")
	cmd.Flags().StringVarP(&options.fileId, "id", "i", "", "id of the file to describe")
	cmd.MarkFlagRequired("id")

	return cmd
}
