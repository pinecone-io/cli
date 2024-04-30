package target

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/spf13/cobra"
)

var targetHelpPart1 string = text.WordWrap(`Many API calls take place in the context of a specific project. 
When using the CLI interactively (i.e. via the device authorization flow) you
should use this command to set the current project context for the CLI.`, 80)

var targetHelpPart2 string = text.WordWrap(pcio.Sprintf(`If you're not sure what values to pass to 
this command, you can discover available projects and organizations by running %s.`, style.Code("pinecone project list")), 80)

var targetHelpPart3 = text.WordWrap(`For automation use cases relying on API-Keys for authentication, there's no need
to specify a project context as the API-Key is already associated with a specific
project in the backend.
`, 80)

var targetHelp = pcio.Sprintf(`%s

%s

%s
`, targetHelpPart1, targetHelpPart2, targetHelpPart3)

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
				log.Info().
					Msg("Outputting target context as table")
				pcio.Printf("To update the context, run %s. The current target context is:\n\n", style.Code("pinecone target --org <org> --project <project>"))
				presenters.PrintTargetContext(state.GetTargetContext())
				return
			}

			orgs, err := dashboard.GetOrganizations()
			if err != nil {
				msg.FailMsg("Failed to get organizations: %s", err)
				exit.Error(err)
			}

			var org dashboard.Organization
			if options.Org != "" {
				// User passed an org flag, need to verify it exists and
				// lookup the id for it.
				org = getOrg(orgs, options.Org)
				if !options.json {
					msg.SuccessMsg("Target org updated to %s", style.Emphasis(org.Name))
				}
				var oldOrg = state.TargetOrg.Get().Name

				// Save the new target org
				state.TargetOrg.Set(&state.TargetOrganization{
					Name: org.Name,
					Id:   org.Id,
				})

				// If the org has changed, reset the project context
				if oldOrg != org.Name {
					state.TargetProj.Set(&state.TargetProject{
						Name: "",
						Id:   "",
					})
				}
			} else {
				// Read the current target org if no org is specified
				// with flags
				org = getOrg(orgs, state.TargetOrg.Get().Name)
			}

			if options.Project != "" {
				// User passed a project flag, need to verify it exists and
				// lookup the id for it.
				proj := getProject(org, options.Project)
				if !options.json {
					msg.SuccessMsg("Target project updated to %s", style.Emphasis(proj.Name))
				}
				state.TargetProj.Set(&state.TargetProject{
					Name: proj.Name,
					Id:   proj.GlobalProject.Id,
				})
			}

			if options.json {
				text.PrettyPrintJSON(state.GetTargetContext())
				return
			}

			pcio.Println()
			presenters.PrintTargetContext(state.GetTargetContext())
		},
	}

	// Required options
	cmd.Flags().StringVarP(&options.Org, "org", "o", "", "Organization name")
	cmd.Flags().StringVarP(&options.Project, "project", "p", "", "Project name")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func getOrg(orgs *dashboard.OrganizationsResponse, orgName string) dashboard.Organization {
	for _, org := range orgs.Organizations {
		if org.Name == orgName {
			return org
		}
	}

	// Join org names for error message
	orgNames := make([]string, len(orgs.Organizations))
	for i, org := range orgs.Organizations {
		orgNames[i] = org.Name
	}

	availableOrgs := strings.Join(orgNames, ", ")
	log.Error().Str("orgName", orgName).Str("avialableOrgs", availableOrgs).Msg("Failed to find organization")
	msg.FailMsg("Failed to find organization %s. Available organizations: %s.", style.Emphasis(orgName), availableOrgs)
	exit.ErrorMsg(pcio.Sprintf("organization %s not found", style.Emphasis(orgName)))
	return dashboard.Organization{}
}

func getProject(org dashboard.Organization, projectName string) dashboard.Project {
	for _, project := range org.Projects {
		if project.Name == projectName {
			return project
		}
	}

	availableProjects := make([]string, len(org.Projects))
	for i, project := range org.Projects {
		availableProjects[i] = project.Name
	}

	msg.FailMsg("Failed to find project %s in org %s. Available projects: %s.", style.Emphasis(projectName), style.Emphasis(org.Name), strings.Join(availableProjects, ", "))
	exit.Error(pcio.Errorf("project %s not found in organization %s", projectName, org.Name))
	return dashboard.Project{}
}
