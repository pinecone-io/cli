package assistant

import (
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type ListKnowledgeModelsCmdOptions struct {
	json bool
}

func NewListAssistantsCmd() *cobra.Command {
	options := ListKnowledgeModelsCmdOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "See the list of assistants in the targeted project",
		GroupID: help.GROUP_ASSISTANT_MANAGEMENT.ID,
		Run: func(cmd *cobra.Command, args []string) {
			modelList, err := assistants.ListAssistants()
			if err != nil {
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(modelList)
				return
			}

			modelCount := len(modelList.Assistants)
			if modelCount == 0 {
				msg.InfoMsg("No assistants found. Create one with %s.\n", style.Code("pinecone assistant create"))
				return
			}

			printTableModels(modelList.Assistants)
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func printTableModels(models []assistants.AssistantModel) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"NAME", "METADATA", "STATUS", "CREATED_AT", "UPDATED_AT"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, model := range models {
		values := []string{model.Name, model.Metadata.ToString(), string(model.Status), model.CreatedAt, model.UpdatedAt}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}
