package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/term"
)

const user = "me"

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to this URL and authorize:\n%v\n", authURL)

	fmt.Print("Enter authorization code: ")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read auth code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to exchange code: %v", err)
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to save token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func main() {
	options := []string{"Check unread mail", "Send email", "Check spam", "Draft mail", "Search"}
	selected := selectOption(options)

	switch options[selected] {
	case "Check unread mail":
		unreadMail(user)
	}
}

func selectOption(options []string) int {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		panic(err)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Println("Select an option:")
	for i, opt := range options {
		fmt.Printf("%s (%d)\n", opt, i+1)
	}
	fmt.Print("Press a number key:")

	for {
		buf := make([]byte, 1)
		_, err := os.Stdin.Read(buf)

		if err != nil {
			panic(err)
		}

		char := buf[0]
		if char == 3 {
			fmt.Println("\nExiting...")
			os.Exit(0)
		}

		if char >= '1' && char <= '9' {
			index := int(char - '1')
			if index < len(options) {
				return index
			}
		}
	}
}
