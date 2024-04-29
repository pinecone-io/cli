package target

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

var targetHelpTemplate string = `Many API calls take place in the context of a specific project. 
When using the CLI interactively (i.e. via the device authorization flow) you
should use this command to set the current project context for the CLI.

If you're not sure what values to pass to this command, you can discover available 
projects and organizations by running %s.

For automation use cases relying on API-Keys for authentication, there's no need
to specify a project context as the API-Key is already associated with a specific
project in the backend.
`
var targetHelp = fmt.Sprintf(targetHelpTemplate, style.Code("pinecone project list"))

type TargetOptions struct {
	Org     string
	Project string
	json    bool
}

func NewTargetCmd() *cobra.Command {
	options := TargetOptions{}

	cmd := &cobra.Command{
		Use:     "target <flags>",
		Short:   "Set context for the CLI",
		GroupID: help.GROUP_START.ID,
		Long:    targetHelp,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug().
				Str("org", options.Org).
				Str("project", options.Project).
				Bool("json", options.json).
				Msg("target command invoked")

			if options.Org == "" && options.Project == "" {
				if options.json {
					log.Info().Msg("Outputting target context as JSON")
					text.PrettyPrintJSON(state.GetTargetContext())
					return
				}
				log.Info().Msg("Outputting target context as table")
				fmt.Printf("To update the context, run %s. The current target context is:\n\n", style.Code("pinecone target --org <org> --project <project>"))
				presenters.PrintTargetContext(state.GetTargetContext())
				return
			}

			orgs, err := dashboard.GetOrganizations()
			if err != nil {
				exit.Error(err)
			}

			var org dashboard.Organization
			if options.Org != "" {
				org, err = getOrg(orgs, options.Org)
				if err != nil {
					exit.Error(err)
				}
				if !options.json {
					fmt.Printf(style.SuccessMsg("Target org updated to %s\n"), style.Emphasis(org.Name))
				}
				state.TargetOrgName.Set(org.Name)
				state.TargetProjectName.Set("")
			} else {
				org, err = getOrg(orgs, state.TargetOrgName.Get())
				if err != nil {
					exit.Error(err)
				}
			}

			if options.Project != "" {
				proj, err := getProject(org, options.Project)
				if err != nil {
					exit.Error(err)
				}
				if !options.json {
					fmt.Printf(style.SuccessMsg("Target project updated to %s\n"), style.Emphasis(proj.Name))
				}
				state.TargetProjectName.Set(proj.Name)
			}

			if options.json {
				text.PrettyPrintJSON(state.GetTargetContext())
				return
			}

			fmt.Println()
			presenters.PrintTargetContext(state.GetTargetContext())
		},
	}

	// Required options
	cmd.Flags().StringVarP(&options.Org, "org", "o", "", "Organization name")
	cmd.Flags().StringVarP(&options.Project, "project", "p", "", "Project name")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func getOrg(orgs *dashboard.OrganizationsResponse, orgName string) (dashboard.Organization, error) {
	for _, org := range orgs.Organizations {
		if org.Name == orgName {
			return org, nil
		}
	}
	return dashboard.Organization{}, fmt.Errorf("organization %s not found", style.Emphasis(orgName))
}

func getProject(org dashboard.Organization, projectName string) (dashboard.Project, error) {
	for _, project := range org.Projects {
		if project.Name == projectName {
			return project, nil
		}
	}
	return dashboard.Project{}, fmt.Errorf("project %s not found in org %s", style.Emphasis(projectName), style.Emphasis(org.Name))
}
