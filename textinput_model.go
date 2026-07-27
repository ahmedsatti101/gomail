package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	input    textinput.Model
	err      error
	quitting bool
}

func textInputModel() string {
	ti := textinput.New()
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 1000
	ti.SetWidth(200)

	m := model{input: ti, err: nil}
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error displaying emails: ", err)
		os.Exit(1)
	}

	m, ok := finalModel.(model)
	if ok {
		return m.input.Value()
	}

	return ""
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "enter":
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	var c *tea.Cursor
	if !m.input.VirtualCursor() {
		c = m.input.Cursor()
		c.Y += lipgloss.Height(m.headerView())
	}

	str := lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.input.View(), m.footerView())
	if m.quitting {
		str += "\n"
	}

	v := tea.NewView(str)
	v.Cursor = c
	return v
}

func (m model) headerView() string {
	return "Supports the same query format as the Gmail search box. For example, 'from:someuser@example.com rfc822msgid: is:unread'\n"
}
func (m model) footerView() string { return "\n(Press Esc or Ctrl+C to quit)" }

