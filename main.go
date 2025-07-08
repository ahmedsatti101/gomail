package main

import (
	"fmt"

	"gomail.com/layout"
)

func main() {
	choice := layout.ChoicesLayout()

	if choice != "" {
		switch choice {
		case "Check unread mail":
			unreadMail("me")
		case "Search mail":
			search()
		case "Send email":
			SendEmail()
		default:
			fmt.Printf("Option '%s' is not available yet\n", choice)
		}
	}
}
