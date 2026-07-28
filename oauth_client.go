package main

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var clientId string
var clientSecret string

func oauthClient(ctx context.Context) *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:9091",
		Scopes:       []string{"https://mail.google.com/", "https://www.googleapis.com/auth/cloud-platform"},
		Endpoint:     google.Endpoint,
	}

	return conf
}
