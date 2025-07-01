package layout

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var choices = []string{"Check unread mail", "Send email", "Check spam", "Draft mail", "Search mail"}

type choicesModel struct {
	cursor int
	choice string
}

func ChoicesLayout() string {
	c := tea.NewProgram(choicesModel{}, tea.WithAltScreen())
	m, err := c.Run()
	
	if err != nil {
		fmt.Printf("Error rendering choices: %v", err)
		os.Exit(1)
	}
	return m.(choicesModel).choice
}

func (m choicesModel) Init() tea.Cmd {
	return nil
}

func (m choicesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}
		case "enter", " ":
			m.choice = choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m choicesModel) View() string {
	s := strings.Builder{}
	s.WriteString("What do you need to do today?\n\n")

	for i := range choices {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(Press q or Ctrl+c to quit)\n")

	return s.String()
}
