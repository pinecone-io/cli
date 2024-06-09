package km

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/knowledge"
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

type KmTargetCmdOptions struct {
	KmName      string
	json        bool
	clear       bool
	interactive bool
}

var kmTargetHelpPart1 string = text.WordWrap(`There are many knowledge model commands which target a specific
knowledge model. This command allows you to set and clear the target knowledge model for performing operations.`, 80)

var targetHelp = pcio.Sprintf("%s\n", kmTargetHelpPart1)

func NewKmTargetCmd() *cobra.Command {
	options := KmTargetCmdOptions{}

	cmd := &cobra.Command{
		Use:     "target <flags>",
		Short:   "Set the target knowledge model",
		Long:    targetHelp,
		GroupID: help.GROUP_KM_TARGETING.ID,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("positional arguments not accepted, please use flags")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug().
				Str("kmName", options.KmName).
				Bool("json", options.json).
				Bool("clear", options.clear).
				Msg("km target command invoked")

			// Clear targets
			if options.clear {
				state.TargetKm.Clear()

				if !options.json {
					msg.SuccessMsg("Target knowledge model cleared.\n")
				}
				printTarget(options.json)
				return
			}

			// Print current target if no knowledge model is specified
			if options.KmName == "" && !options.interactive {
				printTarget(options.json)
				return
			}

			// If model is specified, set target
			modelList, err := knowledge.ListKnowledgeModels()
			if err != nil {
				msg.FailMsg("An error occured while attempting to fetch a list of knowledge models: %s\n", err)
				exit.Error(err)
			}
			if options.KmName != "" {
				// Check if model exists
				modelExists := false
				for _, model := range modelList.KnowledgeModels {
					if model.Name == options.KmName {
						modelExists = true
						break
					}
				}

				if !modelExists {
					availableModels := make([]string, len(modelList.KnowledgeModels))
					for i, model := range modelList.KnowledgeModels {
						availableModels[i] = fmt.Sprintf("'%s'", model.Name)
					}
					sort.Strings(availableModels)
					availableModelsStr := fmt.Sprintf("[%s]", strings.Join(availableModels, ", "))

					msg.FailMsg("Knowledge model %s not found. Available models: %s\n", style.Emphasis(options.KmName), style.Emphasis(availableModelsStr))
					exit.ErrorMsg("knowledge model not found")
					return
				}

				state.TargetKm.Set(&state.TargetKnowledgeModel{Name: options.KmName})

				if !options.json {
					msg.SuccessMsg("Target knowledge model set to %s\n", style.Emphasis(options.KmName))
				}
				printTarget(options.json)
				return
			}

			if options.interactive {
				if len(modelList.KnowledgeModels) == 0 {
					msg.InfoMsg("No knowledge models found. Create one with %s.\n", style.Code("pinecone km create"))
					exit.ErrorMsg("No knowledge models found")
				}

				modelNames := make([]string, len(modelList.KnowledgeModels))
				for i, model := range modelList.KnowledgeModels {
					modelNames[i] = model.Name
				}
				sort.Strings(modelNames)

				selectedModel := uiModelSelector(modelNames)
				if selectedModel == "" {
					// User interrupted selector with ctrl+c
					exit.Success()
				}

				state.TargetKm.Set(&state.TargetKnowledgeModel{Name: selectedModel})
				msg.SuccessMsg("Target knowledge model set to %s\n", style.Emphasis(selectedModel))
				printTarget(options.json)
			} else {
				msg.FailMsg("You must specify a knowledge model with %s or use the %s flag to choose one interactively\n", style.Code("--model"), style.Code("-i"))
				exit.ErrorMsg("no knowledge model specified")
			}
		},
	}

	cmd.Flags().StringVarP(&options.KmName, "model", "m", "", "name of the knowledge model to target")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")
	cmd.Flags().BoolVarP(&options.interactive, "interactive", "i", false, "choose a model interactively")
	cmd.Flags().BoolVar(&options.clear, "clear", false, "clear the target knowledge model")

	return cmd
}

func printTarget(useJson bool) {
	if useJson {
		text.PrettyPrintJSON(state.GetTargetContext())
		return
	}
	presenters.PrintTargetKnowledgeModel(state.GetTargetContext())
}

func uiModelSelector(availableModels []string) string {
	var targetModel string = ""
	prompt := "Choose a model to target"
	listHeight := len(availableModels) + 4
	onQuit := func() {
		pcio.Println("Exiting without targeting a model.")
		pcio.Printf("You can always run %s to change model context later.\n", style.Code("pinecone km target -i"))
	}
	onChoice := func(choice string) string {
		targetModel = choice
		return "Target model: " + choice
	}
	m2 := NewList(availableModels, listHeight, prompt, onQuit, onChoice)
	if _, err := tea.NewProgram(m2).Run(); err != nil {
		pcio.Println("Error selecting knowledge model:", err)
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
