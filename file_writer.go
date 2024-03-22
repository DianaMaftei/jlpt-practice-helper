package main

import (
	"encoding/csv"
	"os"
	"test-go-actions/domain"
)

func writeToCSV(jobPostingsBySource map[string][]domain.JobPosting, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	header := []string{"Title", "Location", "Company", "Field", "Salary", "Experience", "Technologies", "WorkType", "Link", "Source"}
	err = writer.Write(header)
	if err != nil {
		return err
	}

	// Write the job postings
	for source, postings := range jobPostingsBySource {
		for _, job := range postings {
			row := []string{
				job.Title,
				job.Location,
				job.Company,
				job.Field,
				job.Salary,
				job.Experience,
				job.Technologies,
				job.WorkType,
				job.Link,
				source,
			}
			err = writer.Write(row)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
