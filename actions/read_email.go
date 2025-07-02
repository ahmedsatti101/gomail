package actions

import (
	"encoding/base64"
	"fmt"
	"os"

	"gomail.com/utils"
)

func ReadEmail(user, id string) (string, error) {
	srv := utils.CreateService()

	data, dataErr := srv.Users.Messages.Get(user, id).Format("FULL").Do()
	if dataErr != nil {
		return "", fmt.Errorf("Error retriving message info: %v", dataErr)
	}

	var content string

	if data.Payload.Body.Data != "" {
		content = data.Payload.Body.Data
	}

	for _, p1 := range data.Payload.Parts {
		if p1.MimeType == "text/html" || p1.MimeType == "text/plain" {
			content = p1.Body.Data
		} else {
			for _, p2 := range p1.Parts {
				if p2.MimeType == "text/html" || p2.MimeType == "text/plain" {
					content = p2.Body.Data
				}
			}
		}
	}

	message, msgErr := base64.URLEncoding.DecodeString(string(content))
	if msgErr != nil {
		return "", fmt.Errorf("Error decoding email message: %v", msgErr)
	}

	ferr := os.WriteFile("index.html", message, 0644)
	if ferr != nil {
		return "", fmt.Errorf("Error writing to HTML file: %v", ferr)
	}

	return string(message), nil
}
