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

type ListAssistantFilesCmdOptions struct {
	assistant string
	json      bool
}

func NewListAssistantFilesCmd() *cobra.Command {
	options := ListAssistantFilesCmdOptions{}

	cmd := &cobra.Command{
		Use:     "files",
		Short:   "See the list of files in an assistant",
		GroupID: help.GROUP_ASSISTANT_OPERATIONS.ID,
		Run: func(cmd *cobra.Command, args []string) {
			targetAsst := state.TargetAsst.Get().Name
			if targetAsst != "" {
				options.assistant = targetAsst
			}
			if options.assistant == "" {
				msg.FailMsg("You must target an assistant or specify one with the %s flag\n", style.Emphasis("--assistant"))
				exit.Error(fmt.Errorf("no assistant specified"))
			}

			fileList, err := assistants.ListAssistantFiles(options.assistant)
			if err != nil {
				msg.FailMsg("Failed to list files for assistant %s: %s\n", style.Emphasis(options.assistant), err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(fileList)
				pcio.Println(json)
				return
			}

			fileCount := len(fileList.Files)
			if fileCount == 0 {
				msg.InfoMsg("No files found in assistant %s. Add one with %s.\n", style.Emphasis(options.assistant), style.Code("pc assistant file-upload"))
				return
			}

			printTableFiles(fileList.Files)
		},
	}

	cmd.Flags().StringVarP(&options.assistant, "assistant", "a", "", "name of the assistant to list files for")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

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
