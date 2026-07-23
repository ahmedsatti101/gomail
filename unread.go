package main

import (
	"fmt"
	"log"

	"charm.land/bubbles/v2/list"
	"google.golang.org/api/gmail/v1"
)

func unreadMail(service *gmail.Service) {
	fmt.Println("fetching emails...")
	emails, err := service.Users.Messages.List("me").Q("is:unread").MaxResults(500).IncludeSpamTrash(true).Do()
	if err != nil {
		log.Fatalf("Error retriving unread mail: %v", err)
	}

	data := make([]list.Item, 0, len(emails.Messages))

	for _, msgs := range emails.Messages {
		message, err := service.Users.Messages.Get("me", msgs.Id).Format("metadata").Do()
		check(err)

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

		email := email{subject: subject, sender: sender}
		data = append(data, email)
	}

	if len(data) >= 1 {
		List(data)
	} else {
		fmt.Println("No new emails")
	}
}
