package main

import (
	"fmt"
	"encoding/base64"
	"log"
	"os"
	"os/exec"
)

func readEmail(user, id string)  {
	err := os.WriteFile("index.html", []byte(""), 0644)
	if err != nil {
		log.Fatalf("Error creating or writing HTML file: %v", err)
	}

	srv := createService()

	data, dataErr := srv.Users.Messages.Get(user, id).Do()
	if dataErr != nil {
		log.Fatalf("Error retriving message info: %v", dataErr)
	}

	message, msgErr := base64.URLEncoding.DecodeString(data.Payload.Body.Data)
	if msgErr != nil {
		log.Fatalf("Error decoding email message: %v", msgErr)
	}

	ferr := os.WriteFile("index.html", message, 0644)
	if ferr != nil {
		log.Fatalf("Error writing to HTML file: %v", ferr)
	}

	cmd := exec.Command("open", "index.html")
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatalf("Error opening HTML file: %v", err)
	}

	fmt.Println(string(stdout))
}
