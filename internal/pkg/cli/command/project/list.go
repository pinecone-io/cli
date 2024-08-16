package project

import (
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
)

type ListProjectCmdOptions struct {
	json    bool
	all     bool
	orgName string
	orgId   string
}

func NewListProjectsCmd() *cobra.Command {
	options := ListProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list projects in the target org",
		GroupID: help.GROUP_PROJECTS_CRUD.ID,
		Run: func(cmd *cobra.Command, args []string) {
			orgs, err := dashboard.ListOrganizations()
			if err != nil {
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(orgs)
				return
			}

			if options.all {
				printTableAll(orgs)
				return
			}

			if options.orgName != "" {
				for _, org := range orgs.Organizations {
					if org.Name == options.orgName {
						sortProjectsByName(org.Projects)
						printTable(org.Projects)
						return
					}
				}
				exit.Error(pcio.Errorf("organization %s not found", options.orgName))
			}

			if options.orgId != "" {
				for _, org := range orgs.Organizations {
					if org.Id == options.orgId {
						sortProjectsByName(org.Projects)
						printTable(org.Projects)
						return
					}
				}
				exit.Error(pcio.Errorf("organization %s not found", options.orgId))
			}

			targetOrg := state.GetTargetContext().Org
			if targetOrg == "" {
				exit.Error(pcio.Errorf("no target organization set. Please run %s or specify org via flags.", style.Code("pinecone target")))
			}

			for _, org := range orgs.Organizations {
				if org.Name == targetOrg {
					sortProjectsByName(org.Projects)
					printTable(org.Projects)
					return
				}
			}
			// Since the target org is not found, clear the invalid target context
			// to avoid confusion. User can get in this state if they delete the org
			// via some other method (e.g. web, SDK, etc.) and then run this command
			// with saved state that is now stale.
			state.ConfigFile.Clear()
			exit.ErrorMsg(pcio.Sprintf("The target organization %s is not found. Clearing invalid target context. Run %s to see available orgs and %s to set your target context.", style.Emphasis(targetOrg), style.Code("pinecone org list"), style.Code("pinecone target")))
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.orgName, "org_name", "o", "", "name of organization")
	cmd.Flags().StringVarP(&options.orgId, "org_id", "i", "", "id of organization")
	cmd.Flags().BoolVar(&options.all, "all", false, "display projects in all organizations")

	return cmd
}

func sortProjectsByName(projects []dashboard.Project) {
	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})
}

func printTable(projects []dashboard.Project) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"ID", "NAME"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, proj := range projects {
		values := []string{proj.Id, proj.Name}
		pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}

func printTableAll(orgs *dashboard.OrganizationsResponse) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"ORG ID", "ORG NAME", "PROJECT NAME", "PROJECT ID"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	for _, org := range orgs.Organizations {
		for _, proj := range org.Projects {
			values := []string{org.Id, org.Name, proj.Name, proj.Id}
			pcio.Fprintf(writer, strings.Join(values, "\t")+"\n")
		}
	}
	writer.Flush()
}
