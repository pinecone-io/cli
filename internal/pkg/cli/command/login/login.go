package login

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/pinecone-io/cli/internal/pkg/utils/browser"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	pc_oauth2 "github.com/pinecone-io/cli/internal/pkg/utils/oauth2"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Login to Pinecone CLI",
		GroupID: help.GROUP_START.ID,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			da := pc_oauth2.DeviceAuth{}
			authResponse, err := da.GetAuthResponse(ctx)
			if err != nil {
				pcio.Println(err)
				return
			}

			pcio.Printf("Visit %s to authorize the CLI.\n", style.Underline(authResponse.VerificationURIComplete))
			pcio.Println()
			pcio.Printf("The code %s should be displayed on the authorization page.\n", style.HeavyEmphasis(authResponse.UserCode))
			pcio.Println()
			browser.OpenBrowser(authResponse.VerificationURIComplete)

			style.Spinner("Waiting for authorization...", func() error {
				token, err := da.GetDeviceAccessToken(ctx, authResponse)
				if err != nil {
					return err
				}
				secrets.OAuth2Token.Set(token)
				return nil
			})

			pcio.Println()
			accessToken := secrets.OAuth2Token.Get()
			claims, err := pc_oauth2.ParseClaimsUnverified(&accessToken)
			if err != nil {
				log.Error().Msg("Error parsing claims")
				exit.Error(pcio.Errorf("error parsing claims from access token: %s", err))
				return
			}

			pcio.Println(style.SuccessMsg("Logged in as " + style.Emphasis(claims.Email)))

			orgsResponse, err := dashboard.GetOrganizations()
			if err != nil {
				log.Error().Msg("Error fetching organizations")
				exit.Error(pcio.Errorf("error fetching organizations: %s", err))
				return
			}

			targetOrg := postLoginSetTargetOrg(orgsResponse)
			state.TargetOrgName.Set(targetOrg)

			pcio.Println()
			pcio.Printf(style.SuccessMsg("Target org set to %s.\n"), style.Emphasis(targetOrg))

			targetProject := postLoginSetupTargetProject(orgsResponse, targetOrg)
			state.TargetProjectName.Set(targetProject)
			pcio.Printf(style.SuccessMsg("Target project set %s.\n"), style.Emphasis(targetProject))

			pcio.Println()
			pcio.Println(style.CodeHint("Run %s to view or change the target context.", "pinecone target"))

			pcio.Println()
			pcio.Printf("Now try %s to learn about index operations.\n", style.Code("pinecone index -h"))
		},
	}

	return cmd
}

func postLoginSetTargetOrg(orgsResponse *dashboard.OrganizationsResponse) string {
	if len(orgsResponse.Organizations) < 1 {
		log.Debug().Msg("No organizations found. Please create an organization before proceeding.")
		exit.ErrorMsg("No organizations found. Please create an organization before proceeding.")
	}

	// var org dashboard.Organization
	var orgName string
	if len(orgsResponse.Organizations) == 1 {
		orgName = orgsResponse.Organizations[0].Name
		log.Info().Msgf("Only 1 org present so target org set to %s", orgName)
	} else {
		pcio.Println()
		pcio.Println("Many API operations take place in the context of a specific org and project.")
		pcio.Println(pcio.Sprintf("This CLI maintains a piece of state called the %s so it knows which \n", style.Emphasis("target")) +
			"organization and project to use when calling the API on your behalf.")

		orgNames := []string{}
		for _, org := range orgsResponse.Organizations {
			orgNames = append(orgNames, org.Name)
		}
		orgName = uiOrgSelector(orgNames)
	}
	return orgName
}

func postLoginSetupTargetProject(orgs *dashboard.OrganizationsResponse, targetOrg string) string {
	for _, org := range orgs.Organizations {
		if org.Name == targetOrg {
			if len(org.Projects) < 1 {
				log.Debug().Msg("No projects found. Please create a project before proceeding.")
				exit.ErrorMsg("No projects found. Please create a project before proceeding.")
				return ""
			} else if len(org.Projects) == 1 {
				return org.Projects[0].Name
			} else {
				projectItems := []string{}
				for _, proj := range org.Projects {
					projectItems = append(projectItems, proj.Name)
				}
				return uiProjectSelector(projectItems)
			}
		}
	}
	return ""
}

func uiProjectSelector(projectItems []string) string {
	var targetProject string = ""
	m2 := NewList(projectItems, len(projectItems)+6, "Choose a project to target", func() {
		pcio.Println("Exiting without targeting a project.")
		pcio.Printf("You can always run %s to set or change a project context later.\n", style.Code("pinecone target"))
		exit.Success()
	}, func(choice string) string {
		targetProject = choice
		state.TargetProjectName.Set(choice)
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
		pcio.Printf("You can always run %s to set or change a project context later.\n", style.Code("pinecone target"))
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
