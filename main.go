package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const user = "me"

func main() {
	p := tea.NewProgram(choicesModel{}, tea.WithAltScreen())

	m, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	_, ok := m.(choicesModel)

	if ok && m.(choicesModel).choice != "" {
		switch m.(choicesModel).choice {
		case "Check unread mail":
			unreadMail(user)
		default:
			fmt.Println("Not available yet")
		}
	}
}
