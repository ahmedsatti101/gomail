package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

type choicesModel struct {
	cursor int
	choice string
}

var choices = []string{"Check unread mail", "Send email", "Search mail"}

func (ch choicesModel) Init() tea.Cmd {
	return nil
}
func (ch choicesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return ch, tea.Quit

		case "enter":
			ch.choice = choices[ch.cursor]
			return ch, tea.Quit

		case "down", "j":
			ch.cursor++
			if ch.cursor >= len(choices) {
				ch.cursor = 0
			}

		case "up", "k":
			ch.cursor--
			if ch.cursor < 0 {
				ch.cursor = len(choices) - 1
			}
		}
	}

	return ch, nil
}
func (ch choicesModel) View() tea.View {
	s := strings.Builder{}
	s.WriteString("What kind of Bubble Tea would you like to order?\n\n")

	for i := range choices {
		if ch.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return tea.NewView(s.String())
}

func Choices() string {
	p := tea.NewProgram(choicesModel{})
	model, err := p.Run()
	check(err)

	m, ok := model.(choicesModel)
	if ok && m.choice != "" {
		return m.choice
	}

	return ""
}
