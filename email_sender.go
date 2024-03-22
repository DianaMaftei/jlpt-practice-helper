package main

import (
	"bytes"
	"fmt"
	"gopkg.in/mail.v2"
	"os"
	"strconv"
	"strings"
	"test-go-actions/domain"
	"text/template"
	"time"
)

func sendEmail(jobPostings map[string][]domain.JobPosting, csvFilename string, templateName string) error {
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Parse the email template
	funcMap := template.FuncMap{
		"containsJava": func(s string) bool {
			return strings.Contains(s, "Java") && !strings.Contains(s, "JavaScript")
		},
	}
	t, err := template.New(templateName).Funcs(funcMap).ParseFiles(templateName)
	if err != nil {
		return err
	}

	fmt.Println(jobPostings)

	// Prepare the email data
	to := os.Getenv("EMAIL_TO")
	from := os.Getenv("EMAIL_FROM")

	data := struct {
		JobPostings map[string][]domain.JobPosting
	}{
		JobPostings: jobPostings,
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
	m.SetHeader("Subject", "Daily job postings for "+getCurrentDate())
	m.SetBody("text/html", emailBuffer.String())
	m.Attach(csvFilename)

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
