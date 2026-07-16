package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type AuthToken *oauth2.Token

func main() {
	server := &http.Server{
		Addr: ":9091",
	}
	conf := &oauth2.Config{
		ClientID:     "424822125288-2ntgaarra8vaqa4tn15a3kp8oo07ato6.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:9091",
		Scopes:       []string{"https://mail.google.com/"},
		Endpoint:     google.Endpoint,
	}
	consentPageUrl := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	var token AuthToken
	ctx := context.Background()
	srvChan := make(chan struct{}, 1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleAuth(w, r, conf, srvChan)
	})

	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server err: %v", err)
		}
		log.Println("Stopping server...")
	}()

	homeDir, err := os.UserHomeDir()
	check(err)

	_, err = os.Open(filepath.Join(homeDir, ".gomail/creds.json"))

	if errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("Open this link in your browser: %s\n", consentPageUrl)
		<-srvChan
	}

	credsFile, err := os.ReadFile(filepath.Join(homeDir, "/.gomail/creds.json"))
	check(err)

	err = json.Unmarshal(credsFile, &token)
	check(err)

	fmt.Println("fetching your info...")
	if time.Since(token.Expiry) >= time.Hour {
		_, err := updateCreds(conf, token, ctx)
		if err != nil {
			fmt.Println("creds need to be updated")
		}
		check(err)
	}

	googleClient := conf.Client(ctx, token)
	gmailSrv, err := gmail.NewService(ctx, option.WithScopes(gmail.GmailSendScope, gmail.GmailModifyScope, gmail.MailGoogleComScope), option.WithHTTPClient(googleClient))
	check(err)

	profile, err := gmailSrv.Users.GetProfile("me").Do()
	check(err)

	fmt.Printf("profile: %v\n", profile.EmailAddress)
}
