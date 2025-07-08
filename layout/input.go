package layout

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func GetSearchQuery() string {
	// Initialize and run the TUI
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	// Run the program and get the final model
	m, err := p.Run()
	if err != nil {
		log.Fatalf("Error with inputModel: %v", err)
	}

	model := m.(inputModel)

	// Check if user provided a search query. Exit if not
	if len(model.textInput.Value()) == 0 {
		fmt.Println("Please provide a search query")
		os.Exit(1)
	}
	return model.textInput.Value()
}

type errMsg error

type inputModel struct {
	textInput textinput.Model
	err       error
}

func initialModel() inputModel {
	ti := textinput.New()
	ti.Placeholder = "Search"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 100

	return inputModel{
		textInput: ti,
		err:       nil,
	}
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "enter":
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	return fmt.Sprintf(
		"The search query supports the same format as the Gmail search box. (i.e. from:someuser@example.com or is:unread)\n\nYou can combine multiple queries as well. For example, 'category:updates label:unread' which shows unread mail in\nthe updates category\n\n%s\n\n%s",
		m.textInput.View(),
		"(Press Ctrl+c to quit)",
	) + "\n"
}
