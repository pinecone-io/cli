package project

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
)

type ListProjectCmdOptions struct {
	json    bool
	orgName string
	orgId   string
}

func NewListProjectsCmd() *cobra.Command {
	options := ListProjectCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list <command>",
		Short: "list projects in an org",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if options.orgName == "" && options.orgId == "" {
				return fmt.Errorf("organization name or id must be specified")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			orgs, err := dashboard.GetOrganizations(secrets.AccessToken.Get())
			if err != nil {
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(orgs)
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
				exit.Error(fmt.Errorf("organization %s not found", options.orgName))
			}

			if options.orgId != "" {
				for _, org := range orgs.Organizations {
					if org.Id == options.orgId {
						sortProjectsByName(org.Projects)
						printTable(org.Projects)
						return
					}
				}
				exit.Error(fmt.Errorf("organization %s not found", options.orgId))
			}
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().StringVarP(&options.orgName, "org_name", "o", "", "name of organization")
	cmd.Flags().StringVarP(&options.orgId, "org_id", "i", "", "id of organization")

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
	fmt.Fprint(writer, header)

	for _, proj := range projects {
		values := []string{proj.Id, proj.Name}
		fmt.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}
