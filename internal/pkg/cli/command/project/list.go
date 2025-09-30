package project

import (
	"context"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

type ListProjectCmdOptions struct {
	json bool
}

func NewListProjectsCmd() *cobra.Command {
	options := ListProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list all projects in the organization available to the authenticated user",
		GroupID: help.GROUP_PROJECTS.ID,
		Example: help.Examples(`
			pc project list
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ac := sdk.NewPineconeAdminClient()
			ctx := context.Background()

			projects, err := ac.Project.List(ctx)
			if err != nil {
				msg.FailMsg("Failed to list projects: %s\n", err)
				exit.Error(err)
			}

			if options.json {
				json := text.IndentJSON(projects)
				pcio.Println(json)
			} else {
				printTable(projects)
			}
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "Output as JSON")

	return cmd
}

func printTable(projects []*pinecone.Project) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"NAME", "ID", "ORGANIZATION ID", "CREATED AT", "FORCE ENCRYPTION", "MAX PODS"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, proj := range projects {
		values := []string{
			proj.Name,
			proj.Id,
			proj.OrganizationId,
			proj.CreatedAt.String(),
			strconv.FormatBool(proj.ForceEncryptionWithCmek),
			strconv.Itoa(proj.MaxPods)}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}
