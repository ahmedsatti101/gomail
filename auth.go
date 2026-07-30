package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

//go:embed auth-success.html
var f []byte

func handleAuth(w http.ResponseWriter, r *http.Request, conf *oauth2.Config, srvChan chan struct{}) {
	code := r.URL.Query().Get("code")
	if code == "" {
		return
	}

	tokens, err := conf.Exchange(context.Background(), code)
	check(err)

	creds := AuthToken(tokens)
	credsJson, err := json.Marshal(creds)
	check(err)

	homeDir, err := os.UserHomeDir()
	check(err)

	err = os.Chdir(homeDir)
	check(err)

	err = os.Mkdir(".gomail", 0750)
	if !errors.Is(err, os.ErrExist) {
		check(err)
	}

	err = os.WriteFile(filepath.Join(homeDir, ".gomail/creds.json"), credsJson, 0660)
	check(err)

	close(srvChan)

	w.WriteHeader(http.StatusOK)
	w.Write(f)
}
