package main

import (
	"bytes"
	"fmt"
	"gopkg.in/mail.v2"
	"os"
	"strconv"
	"text/template"
	"time"
)

func sendEmail(kanji []Kanji, vocabulary []Vocabulary, grammar []Grammar, videoUrl string, bookUrl string, templateName string) error {
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	t, err := template.New(templateName).ParseFiles(templateName)
	if err != nil {
		return err
	}

	// Prepare the email data
	to := os.Getenv("EMAIL_TO")
	from := os.Getenv("EMAIL_FROM")

	data := struct {
		Kanji      []Kanji
		Vocabulary []Vocabulary
		Grammar    []Grammar
		VideoUrl   string
		BookUrl    string
	}{
		Kanji:      kanji,
		Vocabulary: vocabulary,
		Grammar:    grammar,
		VideoUrl:   videoUrl,
		BookUrl:    bookUrl,
	}

	// Render the email template with the data
	var emailBuffer bytes.Buffer
	err = t.Execute(&emailBuffer, data)
	if err != nil {
		return err
	}

	// Create the mail message
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Daily Japanese Lesson for "+getCurrentDate())
	m.SetBody("text/html", emailBuffer.String())

	// Set up the SMTP client
	smtpPortAsNumber, _ := strconv.Atoi(smtpPort)
	d := mail.NewDialer(smtpHost, smtpPortAsNumber, username, password)

	// Send the email
	err = d.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}

func getCurrentDate() string {
	t := time.Now()

	return fmt.Sprintf("%02d-%02d-%d", t.Day(), t.Month(), t.Year())
}
