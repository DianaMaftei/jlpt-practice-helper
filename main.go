package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"test-go-actions/domain"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	godotenv.Load()

	jobPostingsBySource := make(map[string][]domain.JobPosting)
	jobPostingsBySource = getDummyJobPostings()

	csvFilename := "job_postings.csv"
	err := writeToCSV(jobPostingsBySource, csvFilename)
	if err != nil {
		log.Fatal(err)
	}

	err = sendEmail(jobPostingsBySource, csvFilename, "mail_template.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent with job postings")
}
