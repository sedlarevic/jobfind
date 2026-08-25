package crawler

import (
	"jobfind/model"
	"log"
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

func (c *Crawler) CrawlAll() ([]model.JobPosting, error) {

	var allJobs []model.JobPosting
	for companyName, crawlerFunc := range c.companyRegistry {

		crawledJobs, err := crawlerFunc()

		if err != nil {
			// TODO: Better handling of errors.

			log.Printf("Error crawling %s: %v\n", companyName, err)
			continue
		}

		log.Printf("Job postings for company: %s\n\n", companyName)

		for _, job := range crawledJobs {
			log.Println(job.String() + "\n")
		}

		allJobs = append(allJobs, crawledJobs...)
	}

	return allJobs, nil
}
