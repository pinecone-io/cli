package target

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pinecone-io/cli/internal/pkg/utils/auth"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/login"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/prompt"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
)

var targetHelpPart1 string = text.WordWrap(`Many API calls take place in the context of a specific project. 
When using the CLI interactively (i.e. via the oauth2 authentication flow) you
should use this command to set the current project context for the CLI.`, 80)

var targetHelpPart3 = text.WordWrap(`For automation use cases relying on API-Keys for authentication, there's no need
to specify a project context as the API-Key is already associated with a specific
project in the backend.
`, 80)

var targetHelp = pcio.Sprintf(`%s

%s
`, targetHelpPart1, targetHelpPart3)

type TargetCmdOptions struct {
	Org     string
	Project string
	json    bool
	clear   bool
	show    bool
}

func NewTargetCmd() *cobra.Command {
	options := TargetCmdOptions{}

	cmd := &cobra.Command{
		Use:     "target <flags>",
		Short:   "Set context for the CLI",
		Long:    targetHelp,
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug().
				Str("org", options.Org).
				Str("project", options.Project).
				Bool("json", options.json).
				Msg("target command invoked")

			// Clear targets if clear flag is set
			if options.clear {
				state.ConfigFile.Clear()
				msg.SuccessMsg("Target context has been cleared")
				return
			}

			// Print current target if show flag is set
			if options.show {
				if options.json {
					log.Info().Msg("Outputting target context as JSON")
					json := text.IndentJSON(state.GetTargetContext())
					pcio.Println(json)
					return
				}
				log.Info().
					Msg("Outputting target context as table")
				presenters.PrintTargetContext(state.GetTargetContext())
				return
			}

			// TODO - you don't need to be logged in, you should be able to target with client_id and client_secret
			// But we also need to handle the case where it's a service account rather than a user
			// In that case, we will not have an orgId from the token, we'll need to fetch the org via Admin API
			accessToken, err := auth.Token(cmd.Context())
			if err != nil {
				log.Error().Err(err).Msg("Error retrieving oauth token")
				msg.FailMsg("Error retrieving oauth token: %s", err)
				exit.Error(pcio.Errorf("error retrieving oauth token: %w", err))
				return
			}
			clientId := secrets.ClientId.Get()
			clientSecret := secrets.ClientSecret.Get()
			if accessToken.AccessToken == "" && clientId == "" && clientSecret == "" {
				msg.FailMsg("You must be logged in or have service account credentials configured to set a target context. Run %s to log in, or %s to configure credentials.", style.Code("pc login"), style.Code("pc auth configure"))
				exit.ErrorMsg("You must be logged in or have service account credentials configured to set a target context")
				return
			}
			claims, err := auth.ParseClaimsUnverified(accessToken)
			if err != nil {
				log.Error().Msg("Error parsing claims from access token")
				msg.FailMsg("An auth token was fetched but an error occurred while parsing the token's claims: %s", err)
				exit.Error(pcio.Errorf("error parsing claims from access token: %w", err))
				return
			}
			currentTokenOrgId := claims.OrgId

			ac := sdk.NewPineconeAdminClient()

			// Fetch the user's organizations
			orgs, err := ac.Organization.List(cmd.Context())
			if err != nil {
				log.Error().Msg("Error fetching organizations")
				exit.Error(pcio.Errorf("error fetching organizations: %s", err))
				return
			}

			// Interactive targeting - no options passed
			if options.Org == "" && options.Project == "" && !options.show {

				// Ask the user to choose a target org
				targetOrg := postLoginInteractiveTargetOrg(orgs)
				if targetOrg == nil {
					msg.FailMsg("Failed to target an organization")
					exit.Error(pcio.Errorf("failed to target an organization"))
					return
				} else {
					pcio.Println()
					pcio.Printf(style.SuccessMsg("Target org set to %s.\n"), style.Emphasis(targetOrg.Name))

					// If the org chosen differs from the current orgId in the token, we need to login again
					if currentTokenOrgId != "" && currentTokenOrgId != targetOrg.Id {
						err = login.GetAndSetAccessToken(&targetOrg.Id)
						if err != nil {
							msg.FailMsg("Failed to get access token: %s", err)
							exit.Error(pcio.Errorf("error getting access token: %w", err))
							return
						}
						// Re-create the admin client as the token context has changed
						ac = sdk.NewPineconeAdminClient()
					}
				}

				// Fetch the user's projects
				projects, err := ac.Project.List(cmd.Context())
				if err != nil {
					log.Error().Msg("Error fetching projects")
					exit.Error(pcio.Errorf("error fetching projects: %w", err))
					return
				}

				// Ask the user to choose a target project
				targetProject := postLoginInteractiveTargetProject(projects)
				if targetProject == nil {
					msg.FailMsg("Failed to target a project")
					exit.Error(pcio.Errorf("failed to target a project"))
					return
				} else {
					pcio.Printf(style.SuccessMsg("Target project set %s.\n"), style.Emphasis(targetProject.Name))
					return
				}
			}

			if options.Org != "" {
				// Update the organization target based on passed flag
				var org *pinecone.Organization
				orgs, err := ac.Organization.List(cmd.Context())
				if err != nil {
					msg.FailMsg("Failed to get organizations: %s", err)
					exit.Error(err)
					return
				}

				// User passed an org flag, need to verify it exists and
				// lookup the id for it.
				org, err = getOrg(orgs, options.Org)
				if err != nil {
					msg.FailMsg("Failed to get organization: %s", err)
					exit.Error(err)
					return
				}
				if !options.json {
					msg.SuccessMsg("Target org updated to %s", style.Emphasis(org.Name))
				}
				var oldOrg = state.TargetOrg.Get().Name

				// If the org chosen differs from the current orgId in the token, we need to login again
				if currentTokenOrgId != org.Id {
					err = login.GetAndSetAccessToken(&org.Id)
					if err != nil {
						msg.FailMsg("Failed to get access token: %s", err)
						exit.Error(pcio.Errorf("error getting access token: %w", err))
						return
					}
					// Re-create the admin client as the token context has changed
					ac = sdk.NewPineconeAdminClient()
				}

				// Save the new target org
				state.TargetOrg.Set(state.TargetOrganization{
					Name: org.Name,
					Id:   org.Id,
				})

				// If the org has changed, reset the project context
				if oldOrg != org.Name {
					state.TargetProj.Set(state.TargetProject{
						Name: "",
						Id:   "",
					})
				}
			}

			// Update the project target based on passed flags
			if options.Project != "" {
				// Fetch the user's projects
				projects, err := ac.Project.List(cmd.Context())
				if err != nil {
					log.Error().Msg("Error fetching projects")
					exit.Error(pcio.Errorf("error fetching projects: %w", err))
					return
				}

				// User passed a project flag, need to verify it exists and
				// lookup the id for it.
				proj := getProject(projects, options.Project)
				if !options.json {
					msg.SuccessMsg("Target project updated to %s", style.Emphasis(proj.Name))
				}
				state.TargetProj.Set(state.TargetProject{
					Name: proj.Name,
					Id:   proj.Id,
				})
			}

			// Output JSON if the option was passed
			if options.json {
				json := text.IndentJSON(state.GetTargetContext())
				pcio.Println(json)
				return
			}

			pcio.Println()

			presenters.PrintTargetContext(state.GetTargetContext())
		},
	}

	// Required options
	cmd.Flags().StringVarP(&options.Org, "org", "o", "", "Organization name")
	cmd.Flags().StringVarP(&options.Project, "project", "p", "", "Project name")
	cmd.Flags().BoolVarP(&options.show, "show", "s", false, "Show the current context")
	cmd.Flags().BoolVar(&options.clear, "clear", false, "Clear the target context")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func getOrg(orgs []*pinecone.Organization, orgName string) (*pinecone.Organization, error) {
	for _, org := range orgs {
		if org.Name == orgName {
			return org, nil
		}
	}

	// Join org names for error message
	orgNames := make([]string, len(orgs))
	for i, org := range orgs {
		orgNames[i] = org.Name
	}

	availableOrgs := strings.Join(orgNames, ", ")
	log.Error().Str("orgName", orgName).Str("availableOrgs", availableOrgs).Msg("Failed to find organization")
	msg.FailMsg("Failed to find organization %s. Available organizations: %s.", style.Emphasis(orgName), availableOrgs)
	exit.ErrorMsg(pcio.Sprintf("organization %s not found", style.Emphasis(orgName)))
	return nil, pcio.Errorf("organization %s not found", orgName)
}

func getProject(projects []*pinecone.Project, projectName string) *pinecone.Project {
	for _, project := range projects {
		if project.Name == projectName {
			return project
		}
	}

	availableProjects := make([]string, len(projects))
	for i, project := range projects {
		availableProjects[i] = project.Name
	}
	msg.FailMsg("Failed to find project %s. Available projects: %s.", style.Emphasis(projectName), strings.Join(availableProjects, ", "))
	exit.Error(pcio.Errorf("project %s not found", projectName))
	return nil
}

func postLoginInteractiveTargetOrg(orgsList []*pinecone.Organization) *pinecone.Organization {
	if len(orgsList) < 1 {
		log.Debug().Msg("No organizations found. Please create an organization before proceeding.")
		exit.ErrorMsg("No organizations found. Please create an organization before proceeding.")
	}

	var orgName string
	var organization *pinecone.Organization
	if len(orgsList) == 1 {
		organization = orgsList[0]
		orgName = organization.Name
		log.Info().Msgf("Only 1 organization present. Target organization set to %s", orgName)
	} else {
		pcio.Println("Many API operations take place in the context of a specific org and project.")
		pcio.Println(pcio.Sprintf("This CLI maintains a piece of state called the %s so it knows which \n", style.Emphasis("target")) +
			"organization and project to use when calling the API on your behalf.")

		orgNames := []string{}
		for _, org := range orgsList {
			orgNames = append(orgNames, org.Name)
		}

		orgName = uiOrgSelector(orgNames)
		for _, org := range orgsList {
			if org.Name == orgName {
				state.TargetOrg.Set(state.TargetOrganization{
					Name: org.Name,
					Id:   org.Id,
				})
				organization = org
				break
			}
		}
	}
	return organization
}

func postLoginInteractiveTargetProject(projectList []*pinecone.Project) *pinecone.Project {
	var project *pinecone.Project
	if len(projectList) < 1 {
		log.Debug().Msg("No projects available for organization. Please create a project before proceeding.")
		exit.ErrorMsg("No projects found. Please create a project before proceeding.")
		return nil
	} else if len(projectList) == 1 {
		project = projectList[0]
		state.TargetProj.Set(state.TargetProject{
			Name: projectList[0].Name,
			Id:   projectList[0].Id,
		})
		return projectList[0]
	} else {
		projectItems := []string{}
		for _, proj := range projectList {
			projectItems = append(projectItems, proj.Name)
		}
		projectName := uiProjectSelector(projectItems)

		for _, proj := range projectList {
			if proj.Name == projectName {
				project = proj
				state.TargetProj.Set(state.TargetProject{
					Name: proj.Name,
					Id:   proj.Id,
				})
				return project
			}
		}
	}

	return project
}

func uiProjectSelector(projectItems []string) string {
	var targetProject string = ""
	m2 := prompt.NewList(projectItems, len(projectItems)+6, "Choose a project to target", func() {
		pcio.Println("Exiting without targeting a project.")
		pcio.Printf("You can always run %s to set or change a project context later.\n", style.Code("pc target"))
		exit.Success()
	}, func(choice string) string {
		targetProject = choice
		return "Target project: " + choice
	})
	if _, err := tea.NewProgram(m2).Run(); err != nil {
		pcio.Println("Error running program:", err)
		os.Exit(1)
	}
	return targetProject
}

func uiOrgSelector(orgNames []string) string {
	var orgName string
	m := prompt.NewList(orgNames, len(orgNames)+6, "Choose an organization to target", func() {
		pcio.Println("Exiting without targeting an organization.")
		pcio.Printf("You can always run %s to set or change a project context later.\n", style.Code("pc target"))
		exit.Success()
	}, func(choice string) string {
		orgName = choice
		return "Target organization: " + choice
	})
	if _, err := tea.NewProgram(m).Run(); err != nil {
		pcio.Println("Error running program:", err)
		os.Exit(1)
	}
	return orgName
}
