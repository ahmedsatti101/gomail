package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

func updateCreds(conf *oauth2.Config, oldCreds AuthToken, ctx context.Context) (AuthToken, error) {
	newCreds, err := conf.TokenSource(ctx, oldCreds).Token()

	if err != nil {
		return nil, err
	}

	creds := AuthToken(newCreds)
	credsJson, err := json.Marshal(creds)

	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	err = os.WriteFile(filepath.Join(homeDir, ".gomail/creds.json"), credsJson, 0660)

	if err != nil {
		return nil, err
	}

	return newCreds, nil
}
