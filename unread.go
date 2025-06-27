package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Table style, border style and border color
var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

// Model defining terminal output layout
type unreadMailModel struct {
	table table.Model
}

// Method to use to perform initial I/O. nil means 'no command'
func (m unreadMailModel) Init() tea.Cmd {
	return nil
}

// Updates model when state changes
func (m unreadMailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Selected email with subject %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// Render model
func (m unreadMailModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n(Press q or Ctrl+c to quit)"
}

func unreadMail(user string) {
	service := createService()

	// List unread mail from user's mailinbox
	req, err := service.Users.Messages.List(user).Q("is:unread").MaxResults(10).IncludeSpamTrash(true).Do()
	if err != nil {
		log.Fatalf("Error retriving unread mail: %v", err)
	}

	// Define table columns
	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "Email subject", Width: 70},
		{Title: "Sender", Width: 30},
	}

	// Define table rows
	rows := []table.Row{}

	for _, msgs := range req.Messages {
		// Get each unread email's metadata
		message, err := service.Users.Messages.Get(user, msgs.Id).Format("metadata").Do()
		if err != nil {
			log.Fatalf("%v", err)
		}

		subject := "(No subject)"
		sender := "Unknown"

		if message.Payload != nil {
			for _, header := range message.Payload.Headers {
				switch header.Name {
				case "Subject":
					subject = header.Value
				case "From":
					sender = header.Value
				}
			}
		}

		// Trim sender variable to show only the name of the sender
		if idx := strings.Index(sender, "<"); idx != -1 {
			sender = strings.TrimSpace(sender[:idx])
		}

		// Populate table rows with email id and their subject
		rows = append(rows, []string{msgs.Id, subject, sender})
	}

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

	m := unreadMailModel{t}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error rendering table:", err)
		os.Exit(1)
	}
}
