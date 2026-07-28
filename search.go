package main

import (
	"fmt"
	"log"
	"os"

	"charm.land/bubbles/v2/list"
	"google.golang.org/api/gmail/v1"
)

func search(service *gmail.Service, query string, limit int) {
	fmt.Println("fetching emails...")
	req, err := service.Users.Messages.List("me").Q(query).MaxResults(int64(limit)).Do()
	if err != nil {
		log.Fatalf("Error retriving messages: %v", err)
	}

	if len(req.Messages) == 0 {
		fmt.Println("No emails matching your search. Please try a different search query.")
		os.Exit(1)
	}

	data := make([]list.Item, 0, len(req.Messages))
	for _, msgs := range req.Messages {
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
		fmt.Println("No emails matching your query")
	}
}
