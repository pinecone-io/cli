package project

import (
	"strconv"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

type listProjectCmdOptions struct {
	json bool
}

func NewListProjectsCmd() *cobra.Command {
	options := listProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all projects in the target organization",
		GroupID: help.GROUP_PROJECTS.ID,
		Example: help.Examples(`
			pc project list
		`),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ac := sdk.NewPineconeAdminClient(ctx)

			projects, err := ac.Project.List(ctx)
			if err != nil {
				msg.FailMsg("Failed to list projects: %s\n", err)
				exit.Error(err, "Failed to list projects")
			}

			if options.json {
				json := text.IndentJSON(projects)
				pcio.Println(json)
			} else {
				printTable(projects)
			}
		},
	}

	cmd.Flags().BoolVarP(&options.json, "json", "j", false, "Output as JSON")

	return cmd
}

func printTable(projects []*pinecone.Project) {
	writer := presenters.NewTabWriter()

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
