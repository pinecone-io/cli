package assistant

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/assistants"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type AssistantTargetCmdOptions struct {
	name  string
	clear bool
	show  bool
	json  bool
}

var kmTargetHelpPart1 string = text.WordWrap(`There are many assistant commands which target a specific
assistant. This command allows you to set and clear the target assistant for performing operations.`, 80)

var targetHelp = pcio.Sprintf("%s\n", kmTargetHelpPart1)

func NewAssistantTargetCmd() *cobra.Command {
	options := AssistantTargetCmdOptions{}

	cmd := &cobra.Command{
		Use:     "target <flags>",
		Short:   "Set the target assistant",
		Long:    targetHelp,
		GroupID: help.GROUP_ASSISTANT_TARGETING.ID,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("positional arguments not accepted, please use flags")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug().
				Str("Name", options.name).
				Bool("json", options.json).
				Bool("clear", options.clear).
				Msg("assistant target command invoked")

			// Clear targets
			if options.clear {
				state.TargetAsst.Clear()

				if !options.json {
					msg.SuccessMsg("Target assistant cleared.\n")
				}
				printTarget(options.json)
				return
			}

			// Print current target
			if options.show {
				printTarget(options.json)
				return
			}

			assistantList, err := assistants.ListAssistants()
			if err != nil {
				msg.FailMsg("An error occured while attempting to fetch a list of assistants: %s\n", err)
				exit.Error(err)
			}

			// If an assistant is specified, try to set target
			if options.name != "" {
				assistantExists := false
				for _, model := range assistantList.Assistants {
					if model.Name == options.name {
						assistantExists = true
						break
					}
				}

				if !assistantExists {
					availableAssistants := make([]string, len(assistantList.Assistants))
					for i, model := range assistantList.Assistants {
						availableAssistants[i] = fmt.Sprintf("'%s'", model.Name)
					}
					sort.Strings(availableAssistants)
					availableAssistantsStr := fmt.Sprintf("[%s]", strings.Join(availableAssistants, ", "))

					msg.FailMsg("Assistant %s not found. Available assistants: %s\n", style.Emphasis(options.name), style.Emphasis(availableAssistantsStr))
					exit.ErrorMsg("assistant not found")
					return
				}

				state.TargetAsst.Set(&state.TargetAssistant{Name: options.name})

				msg.SuccessMsg("Target assistant set to %s\n", style.Emphasis(options.name))
				printTarget(options.json)
				return
			}

			// If no assistant is specified and the user has assistants, allow them to select one interactively
			if len(assistantList.Assistants) > 0 {
				assistantNames := make([]string, len(assistantList.Assistants))
				for i, model := range assistantList.Assistants {
					assistantNames[i] = model.Name
				}
				sort.Strings(assistantNames)

				selectedModel := uiModelSelector(assistantNames)
				if selectedModel == "" {
					// User interrupted selector with ctrl+c
					exit.Success()
				}

				state.TargetAsst.Set(&state.TargetAssistant{Name: selectedModel})
				msg.SuccessMsg("Target assistant set to %s\n", style.Emphasis(selectedModel))
				printTarget(options.json)
			} else {
				msg.InfoMsg("No assistants found. Create one with %s.\n", style.Code("pinecone assistant create"))
				exit.ErrorMsg("no assistants found")
			}
		},
	}

	cmd.Flags().StringVarP(&options.name, "name", "n", "", "name of the assistant to target")
	cmd.Flags().BoolVarP(&options.show, "show", "s", false, "show which assistant is currently targeted")
	cmd.Flags().BoolVar(&options.clear, "clear", false, "clear the target assistant")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func printTarget(useJson bool) {
	if useJson {
		text.PrettyPrintJSON(state.GetTargetContext())
		return
	}
	presenters.PrintTargetAssistant(state.GetTargetContext())
}

func uiModelSelector(availableAssistants []string) string {
	var targetModel string = ""
	prompt := "Choose an assistant to target"
	listHeight := len(availableAssistants) + 4
	onQuit := func() {
		pcio.Println("Exiting without targeting an assistant.")
		pcio.Printf("You can always run %s to change assistant context later.\n", style.Code("pinecone assistant target"))
	}
	onChoice := func(choice string) string {
		targetModel = choice
		return "Target assistant: " + choice
	}
	m2 := NewList(availableAssistants, listHeight, prompt, onQuit, onChoice)
	if _, err := tea.NewProgram(m2).Run(); err != nil {
		pcio.Println("Error selecting assistant:", err)
		exit.Error(err)
	}
	return targetModel
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
