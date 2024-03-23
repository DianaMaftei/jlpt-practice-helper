package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	godotenv.Load()

	//jobPostingsBySource := make(map[string][]domain.JobPosting)
	//jobPostingsBySource = getDummyJobPostings()

	Scrape()

	//csvFilename := "job_postings.csv"
	//err := writeToCSV(jobPostingsBySource, csvFilename)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = sendEmail(jobPostingsBySource, csvFilename, "mail_template.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("Email sent with job postings")
}
