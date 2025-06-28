package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gomail.com/layout"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func search() {
	srv := createService()
	q := getSearchQuery()

	fmt.Print("Fetching emails...")
	req, err := srv.Users.Messages.List(user).Q(q).MaxResults(50).Do()
	if err != nil {
		log.Fatalf("Error retriving messages: %v", err)
	}

	if len(req.Messages) == 0 {
		fmt.Println("No available messages. Please try a different query.")
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
		message, err := srv.Users.Messages.Get(user, msgs.Id).Format("metadata").Do()
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

		rows = append(rows, []string{msgs.Id, subject, sender})
	}

	layout.TableLayout(columns, rows)
}

func getSearchQuery() string {
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
		"The search query supports the same format as the Gmail search box. (i.e. from:someuser@example.com or is:unread)\n\n%s\n\n%s",
		m.textInput.View(),
		"(Press Ctrl+c to quit)",
	) + "\n"
}
