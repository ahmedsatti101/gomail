package actions

import (
	"encoding/base64"
	"fmt"
	"os"

	"gomail.com/utils"
)

func ReadEmail(user, id string) (string, error) {
	srv := utils.CreateService()

	data, dataErr := srv.Users.Messages.Get(user, id).Do()
	if dataErr != nil {
		return "", fmt.Errorf("Error retriving message info: %v", dataErr)
	}

	message, msgErr := base64.URLEncoding.DecodeString(data.Payload.Body.Data)
	if msgErr != nil {
		return "", fmt.Errorf("Error decoding email message: %v", msgErr)
	}

	ferr := os.WriteFile("index.html", message, 0644)
	if ferr != nil {
		return "", fmt.Errorf("Error writing to HTML file: %v", ferr)
	}

	return string(message), nil
}
