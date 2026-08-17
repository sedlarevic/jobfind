package crawler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Company struct {
	Name        string       `json:"name"`
	JobPostings []JobPosting `json:"jobpostings"`
}

func (c *Company) String() string {
	var out bytes.Buffer

	out.WriteString(c.Name)
	out.WriteString("\n\n")
	for i, job := range c.JobPostings {
		if i > 0 {
			out.WriteString("\n\n")
		}
		out.WriteString(job.String() + "\n")
	}
	return out.String()
}

type JobPosting struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}

func (jp *JobPosting) String() string {
	return jp.URL + "\n\n" + jp.Text
}

type siteCrawlerFunc func() ([]JobPosting, error)

type Crawler struct {
	registry map[string]siteCrawlerFunc
}

func NewCrawler() *Crawler {
	c := &Crawler{}
	c.registry = map[string]siteCrawlerFunc{
		"NORDEUS": crawlNordeus,
		// Add new companies and their functions
	}
	return c
}

func (c *Crawler) CrawlAllSites(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	for companyName, crawlerFunc := range c.registry {

		jobPostings, err := crawlerFunc()

		company := &Company{Name: companyName, JobPostings: jobPostings}
		if err != nil {
			fmt.Fprintf(w, "Error crawling %s: %v\n", companyName, err)
			continue
		}
		fmt.Fprint(w, company.String())

	}
}

func crawlNordeus() ([]JobPosting, error) {
	var jobPostings []JobPosting
	var siteToVisit = "https://nordeus.com/open-positions/"
	// On every a element which has href attribute call callback

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
		// Print link
		log.Printf("Link found: %q -> %s\n", strings.TrimSpace(e.Text), absoluteURL)
		// Visit link found on page
		collector.Visit(absoluteURL)
	})

	collector.OnHTML(".text-content", func(e *colly.HTMLElement) {
		jobPostings = append(jobPostings, JobPosting{
			URL:  e.Request.URL.String(),
			Text: strings.TrimSpace(e.Text),
		})
	})

	// Before making a request print "Visiting ..."
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
