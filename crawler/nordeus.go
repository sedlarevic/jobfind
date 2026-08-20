package crawler

import (
	"fmt"
	"jobfind/model"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func crawlNordeus() ([]model.JobPosting, error) {
	var jobPostings []model.JobPosting
	var siteToVisit = "https://nordeus.com/open-positions/"

	collector := colly.NewCollector(
		colly.AllowedDomains("www.nordeus.com", "nordeus.com"),
		colly.MaxDepth(1),
	)

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)

		if !strings.Contains(absoluteURL, "/open-positions") || strings.Contains(absoluteURL, "/early_talent_program") || absoluteURL == siteToVisit {
			return
		}

		log.Printf("Link found: %q -> %s\n", strings.TrimSpace(e.Text), absoluteURL)

		collector.Visit(absoluteURL)
	})

	collector.OnHTML(".main", func(e *colly.HTMLElement) {
		if e.Request.URL.String() == siteToVisit {
			return
		}
		jobPostings = append(jobPostings, model.JobPosting{
			Title:       strings.TrimSpace(e.ChildText("h1")),
			URL:         e.Request.URL.String(),
			Description: strings.TrimSpace(e.ChildText(".text-content")),
			CompanyName: "NORDEUS",
			Active:      true,
			FirstSeen:   time.Now(),
			LastSeen:    time.Now(),
		})
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	collector.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	err := collector.Visit(siteToVisit)

	if err != nil {
		return nil, err
	}

	return jobPostings, nil
}
