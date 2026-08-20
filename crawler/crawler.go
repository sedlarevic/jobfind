package crawler

import (
	"fmt"
	"jobfind/model"
	"jobfind/service"
	"net/http"
)

// TODO: Save crawl result to Postgres. New job postings should be added,
// older should just update last_seen column.

// TODO: Every morning at 9.00am crawl all sites to check if there is a difference.
// If yes, alter the table, send notification and display the new result.

// TODO: Load CV PDF (https://github.com/ledongthuc/pdf)

type siteCrawlerFunc func() ([]model.JobPosting, error)

type Crawler struct {
	// NOTE: company name: crawler function
	companyRegistry map[string]siteCrawlerFunc
}

func NewCrawler() *Crawler {
	c := &Crawler{}
	c.companyRegistry = map[string]siteCrawlerFunc{
		"NORDEUS": crawlNordeus,
		// Add new companies and their functions
	}
	return c
}

func (c *Crawler) CrawlAllSites(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	jobPostings := []model.JobPosting{}

	for companyName, crawlerFunc := range c.companyRegistry {

		companyJobPostings, err := crawlerFunc()

		if err != nil {
			fmt.Fprintf(w, "Error crawling %s: %v\n", companyName, err)
			continue
		}

		fmt.Fprintf(w, "Job postings for company: %s\n\n", companyName)

		for _, companyJobPosting := range companyJobPostings {
			fmt.Fprintln(w, companyJobPosting.String()+"\n")
			jobPostings = append(jobPostings, companyJobPosting)
		}

	}

	err := service.SyncJobPostings(jobPostings)

	if err != nil {
		fmt.Fprintf(w, "error syncing job postings")
	}

}
