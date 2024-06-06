package km

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

var helpText = text.WordWrap(`Pinecone Knowledge Engine is a context engine to store and 
retrieve relevant knowledge from millions of documents at scale. 
This API supports creating and managing knowledge models.`, 80)

func NewKmCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "km <command>",
		Short:   "Work with knowledge models",
		Long:    helpText,
		GroupID: help.GROUP_KNOWLEDGE_ENGINE.ID,
	}

	// TODO - add targeting for knowledge model

	cmd.AddCommand(NewCreateKnowledgeModelCmd())
	cmd.AddCommand(NewListKnowledgeModelsCmd())
	cmd.AddCommand(NewDescribeKnowledgeModelCmd())
	cmd.AddCommand(NewDeleteKnowledgeModelCmd())

	// Model-specific commands
	cmd.AddCommand(NewKnowledgeModelChatCmd())
	cmd.AddCommand(NewListKnowledgeFilesCmd())
	cmd.AddCommand(NewDeleteKnowledgeFileCmd())
	cmd.AddCommand(NewUploadKnowledgeFileCmd())
	cmd.AddCommand(NewDescribeKnowledgeFileCmd())

	return cmd
}
