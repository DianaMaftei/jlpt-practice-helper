package main

import (
	"bytes"
	"fmt"
	"gopkg.in/mail.v2"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"
)

func sendEmail(data interface{}, templateName string) error {
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	to := os.Getenv("EMAIL_TO")
	from := os.Getenv("EMAIL_FROM")

	log.Println("Sending email to", to)

	// Use ParseFiles directly
	t, err := template.ParseFiles(templateName)
	if err != nil {
		return err
	}

	// Render the email template with the data
	var emailBuffer bytes.Buffer
	err = t.Execute(&emailBuffer, data)
	if err != nil {
		return err
	}

	// Rest of the function remains the same
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Daily Japanese Lesson for "+getCurrentDate())
	m.SetBody("text/html", emailBuffer.String())

	smtpPortAsNumber, _ := strconv.Atoi(smtpPort)
	d := mail.NewDialer(smtpHost, smtpPortAsNumber, username, password)

	err = d.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Println("Email sent successfully")
	return nil
}

func getCurrentDate() string {
	t := time.Now()

	return fmt.Sprintf("%02d-%02d-%d", t.Day(), t.Month(), t.Year())
}
