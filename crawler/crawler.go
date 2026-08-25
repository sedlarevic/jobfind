package crawler

import (
	"errors"
	"fmt"
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
	/* NOTE: CrawlAll is best-effort.
		Function continues to run if some company crawl error happens.
	 	It returns all succesful crawl results, as well as errors that happened throughout all the crawls.
	*/
	var errs []error
	var allJobs []model.JobPosting

	for companyName, crawlerFunc := range c.companyRegistry {

		crawledJobs, err := crawlerFunc()

		if err != nil {
			crawlErr := fmt.Errorf(
				"crawl %s: %w",
				companyName,
				err,
			)
			log.Println(crawlErr)
			errs = append(errs, crawlErr)
			continue
		}

		allJobs = append(allJobs, crawledJobs...)
	}

	return allJobs, errors.Join(errs...)
}
