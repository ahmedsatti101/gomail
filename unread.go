package main

import (
	"fmt"
	"log"
	"strings"

	"gomail.com/layout"
	"gomail.com/utils"

	"github.com/charmbracelet/bubbles/table"
)

func unreadMail(user string) {
	service := utils.CreateService()

	// Fetch unread mail from user's mailinbox
	fmt.Println("Fetching emails...")
	req, err := service.Users.Messages.List(user).Q("is:unread").MaxResults(50).IncludeSpamTrash(true).Do()
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

		// Populate table rows with email id, their subject & sender info
		rows = append(rows, []string{msgs.Id, subject, sender})
	}

	layout.TableLayout(columns, rows)
}
