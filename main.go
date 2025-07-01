package main

import (
	"fmt"

	"gomail.com/layout"
	"gomail.com/actions"
)

func main() {
	choice := layout.ChoicesLayout()

	if choice != "" {
		switch choice {
		case "Check unread mail":
			actions.UnreadMail("me")
		case "Search mail":
			actions.Search()
		default:
			fmt.Printf("Option '%s' is not available yet", choice)
		}
	}
}
