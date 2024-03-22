package main

import "test-go-actions/domain"

func getDummyJobPostings() map[string][]domain.JobPosting {
	jobPostingsBySource := make(map[string][]domain.JobPosting)

	job1 := domain.JobPosting{
		Title:        "Software Engineer",
		Location:     "New York, NY",
		Company:      "ABC Company",
		Field:        "Software Development",
		Salary:       "80,000 - 100,000",
		Experience:   "3+ years",
		Technologies: "Go, Docker, Kubernetes",
		Link:         "https://www.abccompany.com/jobs/123",
		WorkType:     "Full-time",
		Source:       "Indeed",
	}

	job2 := domain.JobPosting{
		Title:        "DevOps Engineer",
		Location:     "Remote",
		Company:      "XYZ Corporation",
		Field:        "DevOps",
		Experience:   "5+ years",
		Technologies: "Terraform, Ansible, AWS",
		Link:         "https://www.xyzcorp.com/jobs/456",
		WorkType:     "Contract",
		Source:       "Indeed",
	}

	job3 := domain.JobPosting{
		Title:        "DevOps Engineer",
		Location:     "Remote",
		Company:      "XYZ Corporation",
		Field:        "DevOps",
		Salary:       "100,000 - 120,000",
		Technologies: "Terraform, Ansible, Java",
		Link:         "https://www.xyzcorp.com/jobs/456",
		WorkType:     "Contract",
		Source:       "Indeed",
	}

	job4 := domain.JobPosting{
		Title:      "DevOps Engineer",
		Location:   "Remote",
		Company:    "XYZ Corporation",
		Field:      "DevOps",
		Salary:     "100,000 - 120,000",
		Experience: "5+ years",
		Link:       "https://www.xyzcorp.com/jobs/456",
		WorkType:   "Contract",
		Source:     "LinkedIn",
	}

	job5 := domain.JobPosting{
		Title:        "DevOps Engineer",
		Location:     "Remote",
		Company:      "XYZ Corporation",
		Field:        "DevOps",
		Salary:       "100,000 - 120,000",
		Experience:   "5+ years",
		Technologies: "Terraform, Ansible, AWS",
		Link:         "https://www.xyzcorp.com/jobs/456",
		WorkType:     "Contract",
		Source:       "LinkedIn",
	}

	jobPostingsBySource["Indeed"] = append(jobPostingsBySource["Indeed"], job1, job2, job3)
	jobPostingsBySource["LinkedIn"] = append(jobPostingsBySource["LinkedIn"], job4, job5)

	return jobPostingsBySource
}
