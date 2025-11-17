package target

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/login"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/prompt"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v5/pinecone"
	"github.com/spf13/cobra"
)

type targetCmdOptions struct {
	org       string
	orgID     string
	project   string
	projectID string
	json      bool
	clear     bool
	show      bool
}

var (
	targetHelp = help.Long(`
		Set the target organization and project context for the CLI.

		Operations for resources within the control and data plane take place in the context of a specific project.
		After authenticating through the CLI with user login or service account credentials, you can use
		this command to set the target organization or project context for control and data plane operations.

		When using a default API key for authentication, there's no need to specify a project context, because the API 
		key is already associated with a specific organization and project.
	`)

	targetExample = help.Examples(`
		# Interactively target from available organizations and projects
		pc target

		# Target an organization and project by name
		pc target --org "organization-name" -project "project-name"

		# Target a project by name
		pc target --project "project-name"

		# Target an organization and project by ID
		pc target --organization-id "org-id" --project-id "project-id"
	`)
)

func NewTargetCmd() *cobra.Command {
	options := targetCmdOptions{}

	cmd := &cobra.Command{
		Use:     "target",
		Short:   "Set the target organization and project context for the CLI",
		Long:    targetHelp,
		Example: targetExample,
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			log.Debug().
				Str("org", options.org).
				Str("project", options.project).
				Str("organization-id", options.orgID).
				Str("project-id", options.projectID).
				Bool("json", options.json).
				Msg("target command invoked")

			if err := validateTargetOptions(options); err != nil {
				msg.FailMsg("Invalid target options: %s", err)
				exit.Error(err, "Invalid target options")
				return
			}

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
					targetContext := state.GetTargetContext()
					defaultAPIKey := secrets.DefaultAPIKey.Get()
					targetContext.DefaultAPIKey = presenters.MaskHeadTail(defaultAPIKey, 4, 4)
					json := text.IndentJSON(targetContext)
					pcio.Println(json)
					return
				}
				log.Info().
					Msg("Outputting target context as table")

				presenters.PrintTargetContext(state.GetTargetContext())
				return
			}

			// Get the current access token and parse the orgID from the claims
			token, err := oauth.Token(cmd.Context())
			if err != nil {
				msg.FailMsg("Error retrieving oauth token: %s", err)
				exit.Error(err, "Error retrieving oauth token")
			}

			claims, err := oauth.ParseClaimsUnverified(token)
			if err != nil {
				msg.FailMsg("An auth token was fetched but an error occurred while parsing the token's claims: %s", err)
				exit.Error(err, "Error parsing claims from access token")
			}
			currentTokenOrgId := claims.OrgId

			clientId := secrets.ClientId.Get()
			clientSecret := secrets.ClientSecret.Get()
			if token != nil && token.AccessToken == "" && clientId == "" && clientSecret == "" {
				msg.FailMsg("You must be logged in or have service account credentials configured to set a target context. Run %s to log in, or %s to configure credentials.", style.Code("pc login"), style.Code("pc auth configure"))
				exit.ErrorMsg("You must be logged in or have service account credentials configured to set a target context")
			}

			ac := sdk.NewPineconeAdminClient(ctx)

			// Fetch the user's organizations
			orgs, err := ac.Organization.List(cmd.Context())
			if err != nil {
				exit.Error(err, "Error fetching organizations")
			}

			// Interactive targeting - no options passed
			if options.org == "" &&
				options.orgID == "" &&
				options.project == "" &&
				options.projectID == "" {

				// Ask the user to choose a target org
				targetOrg := postLoginInteractiveTargetOrg(orgs)
				if targetOrg == nil {
					msg.FailMsg("Failed to target an organization")
					exit.ErrorMsg("Failed to target an organization")
				} else {
					pcio.Println()
					pcio.Printf(style.SuccessMsg("Target org set to %s.\n"), style.Emphasis(targetOrg.Name))

					// If the org chosen differs from the current orgId in the token, we need to login again
					if currentTokenOrgId != "" && currentTokenOrgId != targetOrg.Id {
						oauth.Logout()
						err = login.GetAndSetAccessToken(ctx, &targetOrg.Id)
						if err != nil {
							msg.FailMsg("Failed to get access token: %s", err)
							exit.Error(err, "Error getting access token")
						}
					}
				}

				ac := sdk.NewPineconeAdminClient(ctx)
				// Fetch the user's projects
				projects, err := ac.Project.List(cmd.Context())
				if err != nil {
					msg.FailMsg("Failed to fetch projects: %s", err)
					exit.Error(err, "error fetching projects")
				}

				// Ask the user to choose a target project
				targetProject := postLoginInteractiveTargetProject(projects)
				if targetProject == nil {
					msg.FailMsg("Failed to target a project")
					exit.ErrorMsg("failed to target a project")
				} else {
					pcio.Printf(style.SuccessMsg("Target project set %s.\n"), style.Emphasis(targetProject.Name))
					return
				}
			}

			// Programmatic targeting - org or orgID flag provided
			if options.org != "" || options.orgID != "" {
				// User organizations were fetched earlier
				var org *pinecone.Organization

				// Use the provided flag to look up the organization
				org, err = getOrgForTarget(orgs, options.org, options.orgID)
				if err != nil {
					msg.FailMsg("Failed to get organization: %s", err)
					exit.Error(err, "Failed to get organization")
				}
				if !options.json {
					msg.SuccessMsg("Target org updated to %s", style.Emphasis(org.Name))
				}
				var oldOrg = state.TargetOrg.Get().Name

				// If the org chosen differs from the current orgId in the token, we need to login again
				if currentTokenOrgId != org.Id {
					oauth.Logout()
					err = login.GetAndSetAccessToken(ctx, &org.Id)
					if err != nil {
						msg.FailMsg("Failed to get access token: %s", err)
						exit.Error(err, "Error getting access token")
					}
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

			// Programmatic targeting - project or projectID flag provided
			if options.project != "" || options.projectID != "" {
				// We need to reinstantiate the admin client to ensure any auth changes that have happened above
				// are properly reflected
				ac := sdk.NewPineconeAdminClient(ctx)

				// Fetch the user's projects
				projects, err := ac.Project.List(cmd.Context())
				if err != nil {
					msg.FailMsg("Error fetching projects: %s", err)
					exit.Error(err, "Error fetching projects")
				}

				// Use the provided flag to look up the project
				project, err := getProjectForTarget(projects, options.project, options.projectID)
				if err != nil {
					msg.FailMsg("Failed to get project: %s", err)
					exit.Error(err, "Failed to get project")
				}
				if !options.json {
					msg.SuccessMsg("Target project updated to %s", style.Emphasis(project.Name))
				}
				state.TargetProj.Set(state.TargetProject{
					Name: project.Name,
					Id:   project.Id,
				})
			}

			// Output JSON if the option was passed
			if options.json {
				targetContext := state.GetTargetContext()
				defaultAPIKey := secrets.DefaultAPIKey.Get()
				targetContext.DefaultAPIKey = presenters.MaskHeadTail(defaultAPIKey, 4, 4)
				json := text.IndentJSON(targetContext)
				pcio.Println(json)
				return
			}

			pcio.Println()

			presenters.PrintTargetContext(state.GetTargetContext())
		},
	}

	// Required options
	cmd.Flags().StringVarP(&options.org, "org", "o", "", "Organization name")
	cmd.Flags().StringVar(&options.orgID, "organization-id", "", "Organization ID")
	cmd.Flags().StringVarP(&options.project, "project", "p", "", "Project name")
	cmd.Flags().StringVar(&options.projectID, "project-id", "", "Project ID")
	cmd.Flags().BoolVarP(&options.show, "show", "s", false, "Show the current context")
	cmd.Flags().BoolVar(&options.clear, "clear", false, "Clear the target context")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func validateTargetOptions(options targetCmdOptions) error {
	// Check organization targeting
	if options.org != "" && options.orgID != "" {
		return pcio.Errorf("cannot specify both --org and --organization-id, use one or the other")
	}

	// Check project targeting
	if options.project != "" && options.projectID != "" {
		return pcio.Errorf("cannot specify both --project and --project-id, use one or the other")
	}

	return nil
}

func getOrgForTarget(orgs []*pinecone.Organization, orgName, orgID string) (*pinecone.Organization, error) {
	var targetOrg *pinecone.Organization
	var searchType string
	var searchValue string

	if orgName != "" {
		// Search by name
		for _, org := range orgs {
			if org.Name == orgName {
				targetOrg = org
				searchValue = orgName
				searchType = "Name"
				break
			}
		}
	} else if orgID != "" {
		// Search by ID
		for _, org := range orgs {
			if org.Id == orgID {
				targetOrg = org
				searchValue = orgID
				searchType = "ID"
				break
			}
		}
	}

	if targetOrg == nil {
		// Join org names for error message
		orgNames := make([]string, len(orgs))
		for i, org := range orgs {
			orgNames[i] = org.Name
		}
		availableOrgs := strings.Join(orgNames, ", ")
		return nil, pcio.Errorf("organization %s: %s not found. Available organizations: %s",
			style.Emphasis(searchType),
			style.Emphasis(searchValue),
			availableOrgs)
	}

	return targetOrg, nil
}

func getProjectForTarget(projects []*pinecone.Project, projectName, projectID string) (*pinecone.Project, error) {
	var targetProject *pinecone.Project
	var searchType string
	var searchValue string

	if projectName != "" {
		// Search by name
		for _, project := range projects {
			if project.Name == projectName {
				targetProject = project
				searchType = "Name"
				searchValue = projectName
				break
			}
		}
	} else if projectID != "" {
		// Search by ID
		for _, project := range projects {
			if project.Id == projectID {
				targetProject = project
				searchType = "ID"
				searchValue = projectID
				break
			}
		}
	}

	if targetProject == nil {
		// Join project names for error message
		projectNames := make([]string, len(projects))
		for i, project := range projects {
			projectNames[i] = project.Name
		}
		availableProjects := strings.Join(projectNames, ", ")
		return nil, pcio.Errorf("project %s: %s not found. Available projects: %s",
			style.Emphasis(searchType),
			style.Emphasis(searchValue),
			availableProjects)
	}

	return targetProject, nil
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
			Name: project.Name,
			Id:   project.Id,
		})
		return project
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
