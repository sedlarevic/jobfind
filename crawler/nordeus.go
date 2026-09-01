package crawler

import (
	"errors"
	"fmt"
	"jobfind/model"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

func crawlNordeus() ([]model.JobPosting, error) {
	var jobPostings []model.JobPosting
	var siteToVisit = "https://nordeus.com/open-positions/"
	var errs []error

	collector := colly.NewCollector(
		colly.AllowedDomains("www.nordeus.com", "nordeus.com"),
		colly.MaxDepth(2),
	)

	// url event
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)

		if !strings.Contains(absoluteURL, "/open-positions") || strings.Contains(absoluteURL, "/early_talent_program") || absoluteURL == siteToVisit {
			return
		}

		log.Printf("Link found: %q -> %s\n", strings.TrimSpace(e.Text), absoluteURL)
		if err := e.Request.Visit(absoluteURL); err != nil {
			log.Printf("could not visit %s: %v", absoluteURL, err)
		}
	})

	// on jobposting page event
	collector.OnHTML(".main", func(e *colly.HTMLElement) {
		if e.Request.URL.String() == siteToVisit {
			return
		}
		jobPostings = append(jobPostings, model.JobPosting{
			Title:       strings.TrimSpace(e.ChildText("h1")),
			URL:         e.Request.URL.String(),
			Description: strings.TrimSpace(e.ChildText(".text-content")),
			CompanyName: "NORDEUS",
		})
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		errs = append(errs, fmt.Errorf("request %s failed: %w", r.Request.URL, err))
	})

	err := collector.Visit(siteToVisit)

	if err != nil {
		errs = append(errs, err)
	}

	return jobPostings, errors.Join(errs...)
}
