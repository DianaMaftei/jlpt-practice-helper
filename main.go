package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"test-go-actions/domain"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jobPostingsBySource := make(map[string][]domain.JobPosting)
	jobPostingsBySource = getDummyJobPostings()

	csvFilename := "job_postings.csv"
	err = writeToCSV(jobPostingsBySource, csvFilename)
	if err != nil {
		log.Fatal(err)
	}

	err = sendEmail(jobPostingsBySource, csvFilename, "mail_template.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent with job postings")
}
