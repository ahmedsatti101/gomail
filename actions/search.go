package actions

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gomail.com/layout"
	"gomail.com/utils"

	"github.com/charmbracelet/bubbles/table"
)

func Search() {
	srv := utils.CreateService()
	q := layout.GetSearchQuery()

	fmt.Println("Fetching emails...")
	req, err := srv.Users.Messages.List("me").Q(q).MaxResults(50).Do()
	if err != nil {
		log.Fatalf("Error retriving messages: %v", err)
	}

	if len(req.Messages) == 0 {
		fmt.Println("No emails matching your search. Please try a different search query.")
		os.Exit(1)
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
		message, err := srv.Users.Messages.Get("me", msgs.Id).Format("metadata").Do()
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
