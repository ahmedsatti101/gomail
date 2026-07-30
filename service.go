package main

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func gmailService(ctx context.Context, token AuthToken, oauth *oauth2.Config) (*gmail.Service, error) {
	googleClient := oauth.Client(ctx, token)
	service, err := gmail.NewService(ctx, option.WithScopes(gmail.GmailSendScope, gmail.GmailModifyScope, gmail.MailGoogleComScope), option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, err
	}

	return service, nil
}
