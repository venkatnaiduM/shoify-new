package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main1() {

	from := "venkatnaidu320@gmail.com"
	password := "dtqj unxa xerl uijk"

	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	to := []string{"venkatnaiduisgod@gmail.com"}

	subject := "Subject: Test Email\n"
	body := "This is a test email sent from Go!"

	auth := smtp.PlainAuth("", from, password, smtpServer)

	message := []byte(subject + "\n" + body)

	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
