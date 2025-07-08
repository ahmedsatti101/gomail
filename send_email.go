package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"gomail.com/layout"
	"gomail.com/utils"
	"google.golang.org/api/gmail/v1"
)

func SendEmail()  {
	layout.SendEmailForm()
	srv := utils.CreateService()
	
	emailValues := layout.GetValues()
	if len(emailValues) == 0 {
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	var message gmail.Message

	headers := make(map[string]string)
	headers["From"] = emailValues[1]
	headers["To"] = emailValues[0]
	headers["Subject"] = emailValues[2]
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	str := ""

	for k, v := range headers {
		str += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	str += "\r\n" + emailValues[3]

	message.Raw = base64.URLEncoding.EncodeToString([]byte(str))

	_, err := srv.Users.Messages.Send("me", &message).Do()

	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	} else {
		fmt.Printf("Email sent to %s\n!", emailValues[0])
	}
}
