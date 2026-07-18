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

	fmt.Println("fetching your info...")
	if time.Since(token.Expiry) >= time.Hour {
		fmt.Println("refreshing sign in details...")
		_, err := updateCreds(oauthClient, token, ctx)
		if err != nil {
			fmt.Printf("Could not update creds: %v", err)
		}
	}

	choice := Choices()
	fmt.Printf("choices: %v\n", choice)
	gmailSrv, err := gmailService(ctx, token, oauthClient)
	check(err)

	profile, err := gmailSrv.Users.GetProfile("me").Do()
	check(err)

	fmt.Printf("profile: %v\n", profile.EmailAddress)
}
