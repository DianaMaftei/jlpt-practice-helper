package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
)

func Scrape() {

	play, err := playwright.Run()
	if err != nil {
		log.Fatal(err)
	}

	browser, err := play.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args: []string{
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		},
	})

	context, err := browser.NewContext()
	if err != nil {
		log.Fatal(err)
	}

	page, err := context.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	_, err = page.Goto("https://www.hays.com.au/job-search/software-engineer-jobs?q=Software%20Engineer")
	if err != nil {
		log.Fatal(err)
	}

	page.WaitForTimeout(2000)

	jobContainers, err := page.Locator(".job-container").All()
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over each selected element and extract job posting data
	for i, element := range jobContainers {
		if i > 3 {
			break
		}

		title, err := element.Locator(".job-descp > h3").TextContent()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(title)
	}

	err = browser.Close()
	if err != nil {
		log.Fatal(err)
	}
}
