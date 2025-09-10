package interactive

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

// ChoiceModel handles choice selection from a list
type ChoiceModel struct {
	prompt       string
	choices      []string
	defaultValue string
	selected     int
	done         bool
	exited       bool
}

// NewChoiceModel creates a new choice selection model
func NewChoiceModel(prompt string, choices []string, defaultValue string) ChoiceModel {
	selected := 0
	// Find default selection
	if defaultValue != "" {
		for i, choice := range choices {
			if choice == defaultValue {
				selected = i
				break
			}
		}
	}
	return ChoiceModel{
		prompt:       prompt,
		choices:      choices,
		defaultValue: defaultValue,
		selected:     selected,
		done:         false,
		exited:       false,
	}
}

func (m ChoiceModel) Init() tea.Cmd {
	return nil
}

func (m ChoiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case msg.String() == "up" || msg.String() == "k":
			if m.selected > 0 {
				m.selected--
			}
		case msg.String() == "down" || msg.String() == "j":
			if m.selected < len(m.choices)-1 {
				m.selected++
			}
		}
	}
	return m, nil
}

func (m ChoiceModel) View() string {
	// Use interactive styles
	questionStyle, _, keyStyle, _ := style.GetInteractiveStyles()

	var s strings.Builder
	s.WriteString(questionStyle.Render(m.prompt))
	s.WriteString("\n")

	if m.done {
		// Show the question and final answer
		selectedChoice := m.choices[m.selected]
		s.WriteString(fmt.Sprintf("  %s\n", style.ResourceName(selectedChoice)))
		return s.String()
	}

	s.WriteString("\n")
	for i, choice := range m.choices {
		cursor := " "
		choiceText := choice
		if m.selected == i {
			cursor = ">"
			choiceText = style.Emphasis(choice)
		}
		s.WriteString(fmt.Sprintf("  %s %s\n", cursor, choiceText))
	}

	s.WriteString(fmt.Sprintf("\nUse %s to select, %s to confirm, %s to exit",
		keyStyle.Render("↑/↓"),
		keyStyle.Render("Enter"),
		keyStyle.Render("Esc/Ctrl+C")))

	return s.String()
}

// GetChoice prompts the user to select from a list of choices
// Returns the selected choice and a boolean indicating if the user wants to exit
func GetChoice(prompt string, choices []string, defaultValue string) (string, bool) {
	m := NewChoiceModel(prompt, choices, defaultValue)

	// Use tea options that properly handle Ctrl+C
	p := tea.NewProgram(m, tea.WithoutSignalHandler())
	finalModel, err := p.Run()
	if err != nil {
		log.Error().Err(err).Msg("Error running choice program")
		return defaultValue, true
	}

	choiceModel, ok := finalModel.(ChoiceModel)
	if !ok {
		log.Error().Msg("Failed to cast final model to ChoiceModel")
		return defaultValue, true
	}

	// Check if user pressed escape, ctrl+c, or q
	if choiceModel.exited {
		return "", true // User wants to exit
	}

	if choiceModel.selected >= 0 && choiceModel.selected < len(choiceModel.choices) {
		return choiceModel.choices[choiceModel.selected], false
	}
	return defaultValue, false
}
