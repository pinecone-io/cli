package interactive

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

// ConfirmationResult represents the result of a confirmation dialog
type ConfirmationResult int

const (
	ConfirmationYes ConfirmationResult = iota
	ConfirmationNo
	ConfirmationQuit
)

// ConfirmationModel handles the user confirmation dialog
type ConfirmationModel struct {
	question string
	choice   ConfirmationResult
	quitting bool
}

// NewConfirmationModel creates a new confirmation dialog model
func NewConfirmationModel(question string) ConfirmationModel {
	return ConfirmationModel{
		question: question,
		choice:   -1, // Invalid state until user makes a choice
	}
}

func (m ConfirmationModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "y", "Y":
			m.choice = ConfirmationYes
			return m, tea.Quit
		case "n", "N":
			m.choice = ConfirmationNo
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ConfirmationModel) View() string {
	if m.quitting {
		return ""
	}
	if m.choice != -1 {
		return ""
	}

	// Use centralized color scheme
	questionStyle, promptStyle, keyStyle, _ := style.GetBrandedConfirmationStyles()

	// Create the confirmation prompt with styled keys
	keys := fmt.Sprintf("%s to confirm, %s to cancel",
		keyStyle.Render("'y'"),
		keyStyle.Render("'n'"))

	return fmt.Sprintf("%s\n%s %s",
		questionStyle.Render(m.question),
		promptStyle.Render("Press"),
		keys)
}

// GetConfirmation prompts the user to confirm an action
// Returns true if the user confirmed with 'y', false if they declined with 'n'
func GetConfirmation(question string) bool {
	m := NewConfirmationModel(question)

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		log.Error().Err(err).Msg("Error running confirmation program")
		return false
	}

	// Get the final model state
	confModel, ok := finalModel.(ConfirmationModel)
	if !ok {
		log.Error().Msg("Failed to cast final model to ConfirmationModel")
		return false
	}

	return confModel.choice == ConfirmationYes
}

// GetConfirmationResult prompts the user to confirm an action and returns the detailed result
// This allows callers to distinguish between "no" and "quit" responses (though both 'n' and 'q' now map to ConfirmationNo)
// Note: Ctrl+C will kill the entire CLI process and is not handled gracefully
func GetConfirmationResult(question string) ConfirmationResult {
	m := NewConfirmationModel(question)

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		log.Error().Err(err).Msg("Error running confirmation program")
		return ConfirmationNo
	}

	// Get the final model state
	confModel, ok := finalModel.(ConfirmationModel)
	if !ok {
		log.Error().Msg("Failed to cast final model to ConfirmationModel")
		return ConfirmationNo
	}

	return confModel.choice
}
