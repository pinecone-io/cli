package interactive

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

// IntInputModel handles integer input with default value
type IntInputModel struct {
	prompt       string
	defaultValue int
	input        string
	done         bool
	exited       bool
}

// NewIntInputModel creates a new integer input model
func NewIntInputModel(prompt string, defaultValue int) IntInputModel {
	return IntInputModel{
		prompt:       prompt,
		defaultValue: defaultValue,
		input:        "",
		done:         false,
		exited:       false,
	}
}

func (m IntInputModel) Init() tea.Cmd {
	return nil
}

func (m IntInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyCtrlC:
			m.done = true
			m.exited = true
			return m, tea.Quit
		case msg.String() == "esc":
			m.done = true
			m.exited = true
			return m, tea.Quit
		case msg.Type == tea.KeyEnter:
			m.done = true
			return m, tea.Quit
		case msg.Type == tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			// Only allow numeric input
			keypress := msg.String()
			if len(keypress) == 1 && (keypress >= "0" && keypress <= "9") {
				m.input += keypress
			}
		}
	}
	return m, nil
}

func (m IntInputModel) View() string {
	// Use interactive styles
	questionStyle, inputStyle, _, _ := style.GetInteractiveStyles()

	prompt := fmt.Sprintf("%s [%d]: ", m.prompt, m.defaultValue)

	if m.done {
		// Show the question and final answer
		finalAnswer := m.input
		if finalAnswer == "" {
			finalAnswer = fmt.Sprintf("%d", m.defaultValue)
		}
		return fmt.Sprintf("%s%s\n",
			questionStyle.Render(prompt),
			style.ResourceName(finalAnswer))
	}

	return fmt.Sprintf("%s%s_",
		questionStyle.Render(prompt),
		inputStyle.Render(m.input))
}

// GetIntInput prompts the user for integer input with a default value
// Returns the input value and a boolean indicating if the user wants to exit
func GetIntInput(prompt string, defaultValue int) (int, bool) {
	m := NewIntInputModel(prompt, defaultValue)

	// Use tea options that properly handle Ctrl+C
	p := tea.NewProgram(m, tea.WithoutSignalHandler())
	finalModel, err := p.Run()
	if err != nil {
		log.Error().Err(err).Msg("Error running int input program")
		return defaultValue, true
	}

	intModel, ok := finalModel.(IntInputModel)
	if !ok {
		log.Error().Msg("Failed to cast final model to IntInputModel")
		return defaultValue, true
	}

	// Check if user pressed escape, ctrl+c, or q
	if intModel.exited {
		return 0, true // User wants to exit
	}

	input := strings.TrimSpace(intModel.input)
	if input == "" {
		return defaultValue, false
	}

	value, err := strconv.Atoi(input)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing integer input")
		return defaultValue, false
	}

	return value, false
}
