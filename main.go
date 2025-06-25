package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const user = "me"

var choices = []string{"Check unread mail", "Send email", "Check spam", "Draft mail", "Search mail"}

type model struct {
	cursor int
	choice string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m model) View() string {
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
	s.WriteString("\n(Press q to quit)\n")

	return s.String()
}

func main() {
	p := tea.NewProgram(model{})

	m, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	_, ok := m.(model)

	if ok && m.(model).choice != "" {
		switch m.(model).choice {
		case "Check unread mail":
			unreadMail(user)
		case "Send email":
			fmt.Println("Not available yet")
		case "Check spam":
			fmt.Println("Not available yet")
		case "Draft mail":
			fmt.Println("Not available yet")
		case "Search mail":
			fmt.Println("Not available yet")
		}
	}
}
