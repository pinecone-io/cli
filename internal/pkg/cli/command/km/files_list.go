package km

import (
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/knowledge"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type ListKnowledgeFilesCmdOptions struct {
	json   bool
	kmName string
}

func NewListKnowledgeFilesCmd() *cobra.Command {
	options := ListKnowledgeFilesCmdOptions{}

	cmd := &cobra.Command{
		Use:     "files",
		Short:   "See the list of files in a knowledge model",
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

			fileList, err := knowledge.ListKnowledgeModelFiles(options.kmName)
			if err != nil {
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(fileList)
				return
			}
			pcio.Printf("Found %d files in knowledge model %s\n", len(fileList.Files), options.kmName)
			printTableFiles(fileList.Files)
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.kmName, "name", "n", "", "name of the knowledge model to list files for")

	return cmd
}

func printTableFiles(files []knowledge.KnowledgeFileModel) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"NAME", "ID", "METADATA", "CREATED_ON", "UPDATED_ON", "STATUS"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, file := range files {
		values := []string{
			file.Name,
			file.Id,
			file.Metadata.ToString(),
			file.CreatedOn,
			file.UpdatedOn,
			string(file.Status),
		}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}
