package main

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func createService() *gmail.Service {
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

	return srv
}
