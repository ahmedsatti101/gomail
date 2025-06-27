package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func search() {
	srv := createService()
	q := getSearchQuery()

	req, err := srv.Users.Messages.List(user).Q(q).Do()
	if err != nil {
		log.Fatalf("Error retriving messages: %v", err)
	}

	if len(req.Messages) == 0 {
		fmt.Println("No available messages. Please try a different query.")
	}

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
		
		fmt.Println(msgs.Id)
		fmt.Println(subject)
		fmt.Println(sender)
	}
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
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEnter:
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
