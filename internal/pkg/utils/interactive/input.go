package interactive

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

// InputModel handles text input with default value
type InputModel struct {
	prompt       string
	defaultValue string
	input        string
	done         bool
	exited       bool
}

// NewInputModel creates a new text input model
func NewInputModel(prompt, defaultValue string) InputModel {
	return InputModel{
		prompt:       prompt,
		defaultValue: defaultValue,
		input:        "",
		done:         false,
		exited:       false,
	}
}

func (m InputModel) Init() tea.Cmd {
	return nil
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}
	}
	return m, nil
}

func (m InputModel) View() string {
	// Use interactive styles
	questionStyle, _, _, _ := style.GetInteractiveStyles()

	prompt := m.prompt
	if m.defaultValue != "" {
		prompt += fmt.Sprintf(" [%s]", style.Emphasis(m.defaultValue))
	}
	prompt += ": "

	if m.done {
		// Show the question and final answer
		finalAnswer := m.input
		if finalAnswer == "" {
			finalAnswer = m.defaultValue
		}
		return fmt.Sprintf("%s%s\n",
			questionStyle.Render(prompt),
			style.ResourceName(finalAnswer))
	}

	return fmt.Sprintf("%s%s_",
		questionStyle.Render(prompt),
		m.input)
}

// GetInput prompts the user for text input with a default value
// Returns the input value and a boolean indicating if the user wants to exit
func GetInput(prompt, defaultValue string) (string, bool) {
	m := NewInputModel(prompt, defaultValue)

	// Use tea options that properly handle Ctrl+C
	p := tea.NewProgram(m, tea.WithoutSignalHandler())
	finalModel, err := p.Run()
	if err != nil {
		log.Error().Err(err).Msg("Error running input program")
		return defaultValue, true
	}

	inputModel, ok := finalModel.(InputModel)
	if !ok {
		log.Error().Msg("Failed to cast final model to InputModel")
		return defaultValue, true
	}

	// Check if user pressed escape, ctrl+c, or q
	if inputModel.exited {
		return "", true // User wants to exit
	}

	input := strings.TrimSpace(inputModel.input)
	if input == "" {
		return defaultValue, false
	}
	return input, false
}
