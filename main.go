package main

import (
	"fmt"

	"gomail.com/layout"
)

const user = "me"

func main() {
	choice := layout.ChoicesLayout()

	if choice != "" {
		switch choice {
		case "Check unread mail":
			unreadMail(user)
		case "Search mail":
			search()
		case "Read mail":
			readEmail(user, "197bdf4f1412b838")
		default:
			fmt.Printf("Option '%s' is not available yet", choice)
		}
	}
}
