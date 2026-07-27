package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
)

type AuthToken *oauth2.Token

func main() {
	ctx := context.Background()
	server := &http.Server{
		Addr: ":9091",
	}
	oauthClient := oauthClient(ctx)
	var token AuthToken
	srvChan := make(chan struct{}, 1)

	limitFlag := flag.Int("limit", 50, "Specify the limit flag to limit how many emails are retrieved. Max value is 500.")
	flag.Parse()
	if *limitFlag > 500 {
		fmt.Println("The maximum number of emails that can be retrieved is 500.")
		os.Exit(1)
	}

	consentPageUrl := oauthClient.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleAuth(w, r, oauthClient, srvChan)
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

	if time.Since(token.Expiry) >= time.Hour {
		fmt.Println("refreshing sign in details...")
		_, err := updateCreds(oauthClient, token, ctx)
		if err != nil {
			fmt.Printf("Could not update creds: %v", err)
		}
	}

	gmailSrv, err := gmailService(ctx, token, oauthClient)
	check(err)

	choice := Choices()
	switch choice {
	case "Check unread mail":
		unreadMail(gmailSrv, *limitFlag)
	case "Search mail":
		q := textInputModel()
		if q != "" {
			search(gmailSrv, q, *limitFlag)
		}
	default:
		fmt.Println("Not implemented")
	}
}
