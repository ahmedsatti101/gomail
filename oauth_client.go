package main

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func oauthClient(ctx context.Context) *oauth2.Config {
	clientId, err := secretManager(ctx, "projects/424822125288/secrets/gomail-client-id/versions/latest")
	check(err)
	clientSecret, err := secretManager(ctx, "projects/424822125288/secrets/gomail-client-secret/versions/latest")
	check(err)

	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:9091",
		Scopes:       []string{"https://mail.google.com/", "https://www.googleapis.com/auth/cloud-platform"},
		Endpoint:     google.Endpoint,
	}

	return conf
}
