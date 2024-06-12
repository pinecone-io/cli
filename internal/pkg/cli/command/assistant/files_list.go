package assistant

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

type ListKnowledgeFilesCmdOptions struct {
	json bool
	name string
}

func NewListAssistantFilesCmd() *cobra.Command {
	options := ListKnowledgeFilesCmdOptions{}

	cmd := &cobra.Command{
		Use:     "files",
		Short:   "See the list of files in an assistant",
		GroupID: help.GROUP_ASSISTANT_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			targetKm := state.TargetAsst.Get().Name
			if targetKm != "" {
				options.name = targetKm
			}
			if options.name == "" {
				msg.FailMsg("You must target an assistant or specify one with the %s flag\n", style.Emphasis("--model"))
				exit.Error(fmt.Errorf("no assistant specified"))
			}

			fileList, err := assistants.ListAssistantFiles(options.name)
			if err != nil {
				msg.FailMsg("Failed to list files for assistant %s: %s\n", style.Emphasis(options.name), err)
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(fileList)
				return
			}

			fileCount := len(fileList.Files)
			if fileCount == 0 {
				msg.InfoMsg("No files found in assistant %s. Add one with %s.\n", style.Emphasis(options.name), style.Code("pinecone assistant file-upload"))
				return
			}

			printTableFiles(fileList.Files)
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.name, "model", "m", "", "name of the assistant to list files for")

	return cmd
}

func printTableFiles(files []assistants.AssistantFileModel) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"NAME", "ID", "METADATA", "CREATED_ON", "UPDATED_ON", "STATUS", "SIZE"}
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
			fmt.Sprintf("%d", file.Size),
		}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}
