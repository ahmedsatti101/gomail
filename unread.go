package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func unreadMail(user string)  {
	ctx := context.Background()
	f, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials: %v", err)
	}

	config, err := google.ConfigFromJSON(f, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse credentials file: %v", err)
	}

	client := getClient(config)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		log.Fatalf("Error initialising gmail service %v", err)
	}

	r, err := srv.Users.Messages.List(user).Q("is:unread").MaxResults(10).IncludeSpamTrash(true).Do()
	if err != nil {
		log.Fatalf("Error retriving unread mail: %v", err)
	}

	messages := [][]*gmail.MessagePartHeader{}
	
	for _, i := range r.Messages {
		message, err := srv.Users.Messages.Get(user, i.Id).Format("full").Do()
		if err != nil {
			log.Fatalf("%v", err)
		}
		messages = append(messages, message.Payload.Headers)
	}

	for _, i := range messages {
		for _, i := range i {
			if i.Name == "Subject" {
				fmt.Println(i.Value)
			}
		}
	}
}
