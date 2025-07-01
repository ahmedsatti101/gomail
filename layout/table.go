package layout

import (
	"fmt"
	"os"
	"os/exec"

	"gomail.com/actions"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TableLayout(columns []table.Column, rows []table.Row) string {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#00000000")).
		Background(lipgloss.Color("#ffffff")).
		Bold(false)
	t.SetStyles(s)

	m := tableModel{t}
	_, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error rendering table:", err)
		os.Exit(1)
	}
	return m.table.SelectedRow()[0]
}

// Table style, border style and border color
var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

// Model defining terminal output layout
type tableModel struct {
	table table.Model
}

// Method to use to perform initial I/O. nil means 'no command'
func (m tableModel) Init() tea.Cmd {
	return nil
}

// Updates model when state changes
func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			selectedId := m.table.SelectedRow()[0]
			return m, fetchEmail("me", selectedId)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// Render model
func (m tableModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n(Press q or Ctrl+c to quit)"
}

type emailMsg struct {
	msg string
	err error
}

func fetchEmail(user, id string) tea.Cmd {
	return func() tea.Msg {
		msg, err := actions.ReadEmail(user, id)
		if err != nil {
			return emailMsg{msg, err}
		}

		if len(msg) > 0 {
			cmd := exec.Command("open", "index.html")
			if err := cmd.Start(); err != nil {
				return emailMsg{msg, err}
			}
		}
		return emailMsg{msg, nil}
	}
}
