package km

import (
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type ListKnowledgeModelsCmdOptions struct {
	json bool
}

func NewListKnowledgeModelsCmd() *cobra.Command {
	options := ListKnowledgeModelsCmdOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "See the list of knowledge models in the targeted project",
		GroupID: help.GROUP_KM_MANAGEMENT.ID,
		Run: func(cmd *cobra.Command, args []string) {
			modelList, err := knowledge.ListKnowledgeModels()
			if err != nil {
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(modelList)
				return
			}

			modelCount := len(modelList.KnowledgeModels)
			if modelCount == 0 {
				pcio.Println("No knowledge models found")
				return
			}

			printTableModels(modelList.KnowledgeModels)
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func printTableModels(models []knowledge.KnowledgeModel) {
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
