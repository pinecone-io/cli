package target

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/login"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	pc_oauth2 "github.com/pinecone-io/cli/internal/pkg/utils/oauth2"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
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

			access_token := secrets.OAuth2Token.Get()
			if access_token.AccessToken == "" {
				msg.FailMsg("You must be logged in to set a target context. Run %s to log in.", style.Code("pc login"))
				exit.ErrorMsg("You must be logged in to set a target context")
				return
			}
			claims, err := pc_oauth2.ParseClaimsUnverified(&access_token)
			if err != nil {
				log.Error().Msg("Error parsing claims from access token")
				msg.FailMsg("An auth token was fetched but an error occurred while parsing the token's claims: %s", err)
				exit.Error(pcio.Errorf("error parsing claims from access token: %w", err))
				return
			}
			currentTokenOrgId := claims.OrgId

			// Interactive targeting if logged in
			if options.Org == "" && options.Project == "" && !options.show {
				// Fetch the user's organizations and projects
				orgsResponse, err := dashboard.ListOrganizations()
				if err != nil {
					log.Error().Msg("Error fetching organizations")
					exit.Error(pcio.Errorf("error fetching organizations: %s", err))
					return
				}

				// Ask the user to choose a target org
				targetOrg := postLoginSetTargetOrg(orgsResponse)
				pcio.Println()
				pcio.Printf(style.SuccessMsg("Target org set to %s.\n"), style.Emphasis(targetOrg.Name))

				// If the org chosen differs from the current orgId in the token, we need to login again
				if currentTokenOrgId != targetOrg.Id {
					err = login.GetAndSetAccessToken(&targetOrg.Id)
					if err != nil {
						msg.FailMsg("Failed to get access token: %s", err)
						exit.Error(pcio.Errorf("error getting access token: %w", err))
						return
					}
				}

				// Fetch orgs again since this currently also populates projects
				// TODO: uncouple the org and project fetching
				orgsResponse, err = dashboard.ListOrganizations()
				if err != nil {
					log.Error().Msg("Error fetching organizations")
					exit.Error(pcio.Errorf("error fetching organizations: %s", err))
					return
				}

				// Ask the user to choose a target project
				targetProject := postLoginSetupTargetProject(orgsResponse, targetOrg.Name)
				pcio.Printf(style.SuccessMsg("Target project set %s.\n"), style.Emphasis(targetProject))
				return
			}

			orgs, err := dashboard.ListOrganizations()
			if err != nil {
				msg.FailMsg("Failed to get organizations: %s", err)
				exit.Error(err)
				return
			}

			// Update the organization target
			var org dashboard.Organization
			if options.Org != "" {
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

				// If the org has changed, we need to swap to a new access token
				// If the org chosen differs from the current orgId in the token, we need to login again
				if currentTokenOrgId != org.Id {
					pcio.Printf("GET/SET NEW ACCESS TOKEN FOR ORG %s WITH ID %s\n", org.Name, org.Id)
					err = login.GetAndSetAccessToken(&org.Id)
					if err != nil {
						msg.FailMsg("Failed to get access token: %s", err)
						exit.Error(pcio.Errorf("error getting access token: %w", err))
						return
					}
				}

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
				org, err = getOrg(orgs, state.TargetOrg.Get().Name)
				if err != nil {
					msg.FailMsg("Failed to get organization: %s", err)
					exit.Error(err)
					return
				}
			}

			// Update the project target
			if options.Project != "" {
				// User passed a project flag, need to verify it exists and
				// lookup the id for it.
				proj := getProject(org, options.Project)
				if !options.json {
					msg.SuccessMsg("Target project updated to %s", style.Emphasis(proj.Name))
				}
				state.TargetProj.Set(&state.TargetProject{
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

func getOrg(orgs *dashboard.OrganizationsResponse, orgName string) (dashboard.Organization, error) {
	for _, org := range orgs.Organizations {
		if org.Name == orgName {
			return org, nil
		}
	}

	// Join org names for error message
	orgNames := make([]string, len(orgs.Organizations))
	for i, org := range orgs.Organizations {
		orgNames[i] = org.Name
	}

	availableOrgs := strings.Join(orgNames, ", ")
	log.Error().Str("orgName", orgName).Str("availableOrgs", availableOrgs).Msg("Failed to find organization")
	msg.FailMsg("Failed to find organization %s. Available organizations: %s.", style.Emphasis(orgName), availableOrgs)
	exit.ErrorMsg(pcio.Sprintf("organization %s not found", style.Emphasis(orgName)))
	return dashboard.Organization{}, pcio.Errorf("organization %s not found", orgName)
}

func getProject(org dashboard.Organization, projectName string) dashboard.Project {
	if org.Projects != nil {
		orgProjects := *org.Projects
		for _, project := range orgProjects {
			if project.Name == projectName {
				return project
			}
		}

		availableProjects := make([]string, len(orgProjects))
		for i, project := range orgProjects {
			availableProjects[i] = project.Name
		}
		msg.FailMsg("Failed to find project %s in org %s. Available projects: %s.", style.Emphasis(projectName), style.Emphasis(org.Name), strings.Join(availableProjects, ", "))
		exit.Error(pcio.Errorf("project %s not found in organization %s", projectName, org.Name))
		return dashboard.Project{}
	} else {
		msg.FailMsg("Cannot load projects for current organization. Please log into organization to retrieve projects.")
		exit.Error(pcio.Errorf("project %s not found in organization %s", projectName, org.Name))
		return dashboard.Project{}
	}
}

func postLoginSetTargetOrg(orgsResponse *dashboard.OrganizationsResponse) dashboard.Organization {
	if len(orgsResponse.Organizations) < 1 {
		log.Debug().Msg("No organizations found. Please create an organization before proceeding.")
		exit.ErrorMsg("No organizations found. Please create an organization before proceeding.")
	}

	// var org dashboard.Organization
	var orgName string
	var organization dashboard.Organization
	if len(orgsResponse.Organizations) == 1 {
		orgName = orgsResponse.Organizations[0].Name
		log.Info().Msgf("Only 1 org present so target org set to %s", orgName)
	} else {
		pcio.Println("Many API operations take place in the context of a specific org and project.")
		pcio.Println(pcio.Sprintf("This CLI maintains a piece of state called the %s so it knows which \n", style.Emphasis("target")) +
			"organization and project to use when calling the API on your behalf.")

		orgNames := []string{}
		for _, org := range orgsResponse.Organizations {
			orgNames = append(orgNames, org.Name)
		}

		orgName = uiOrgSelector(orgNames)
		for _, org := range orgsResponse.Organizations {
			if org.Name == orgName {
				state.TargetOrg.Set(&state.TargetOrganization{
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

func postLoginSetupTargetProject(orgs *dashboard.OrganizationsResponse, targetOrg string) string {
	for _, org := range orgs.Organizations {
		if org.Name == targetOrg {
			// Attempt to fetch projects for a given org. We assume the token retrieved from the
			// device auth flow is valid for listing projects for the given org.
			projectsList, err := dashboard.ListProjects(org.Id)
			if err != nil {
				log.Error().Msg("Error fetching projects")
				exit.Error(pcio.Errorf("error fetching projects: %w", err))
				return ""
			}

			orgProjects := projectsList.Projects
			if len(orgProjects) < 1 {
				log.Debug().Msg("No projects available for organization. Please create a project before proceeding.")
				exit.ErrorMsg("No projects found. Please create a project before proceeding.")
				return ""
			} else {
				if len(orgProjects) < 1 {
					log.Debug().Msg("No projects found. Please create a project before proceeding.")
					exit.ErrorMsg("No projects found. Please create a project before proceeding.")
					return ""
				} else if len(orgProjects) == 1 {
					state.TargetProj.Set(&state.TargetProject{
						Name: orgProjects[0].Name,
						Id:   orgProjects[0].Id,
					})
					return orgProjects[0].Name
				} else {
					projectItems := []string{}
					for _, proj := range orgProjects {
						projectItems = append(projectItems, proj.Name)
					}
					projectName := uiProjectSelector(projectItems)

					for _, proj := range orgProjects {
						if proj.Name == projectName {
							state.TargetProj.Set(&state.TargetProject{
								Name: proj.Name,
								Id:   proj.Id,
							})
							return proj.Name
						}
					}
				}
			}
		}
	}
	return ""
}

func uiProjectSelector(projectItems []string) string {
	var targetProject string = ""
	m2 := NewList(projectItems, len(projectItems)+6, "Choose a project to target", func() {
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
	m := NewList(orgNames, len(orgNames)+6, "Choose an organization to target", func() {
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

type ListModel struct {
	list     list.Model
	choice   string
	quitting bool
	onQuit   func()
	onChoice func(string) string
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			m.onChoice(m.choice)
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	if m.choice != "" {
		return "sounds great " + m.choice
	}
	if m.quitting {
		m.onQuit()
		return ""
	}
	return "\n" + m.list.View()
}

func mapStringsToItems(strings []string) []list.Item {
	items := make([]list.Item, len(strings))
	for i, s := range strings {
		items[i] = item(s)
	}
	return items
}

func NewList(items []string, listHeight int, title string, onQuit func(), onChoice func(string) string) ListModel {
	const defaultWidth = 20

	l := list.New(mapStringsToItems(items), itemDelegate{}, defaultWidth, listHeight)
	l.SetShowHelp(false)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle

	return ListModel{
		list:     l,
		onQuit:   onQuit,
		onChoice: onChoice,
	}
}

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(0)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(3)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("5"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := pcio.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	pcio.Fprint(w, fn(str))
}
