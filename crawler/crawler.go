package crawler

import (
	"fmt"
	"jobfind/crawler/companies"
	"jobfind/model"
	"log/slog"
	"time"
)

type siteCrawlerFunc func() ([]model.JobPosting, error)

type CompanyCrawlResult struct {
	Company     string
	JobPostings []model.JobPosting
	Err         error
}

type Crawler struct {
	// NOTE: company name: crawler function
	companyRegistry map[string]siteCrawlerFunc
}

func NewCrawler() *Crawler {
	c := &Crawler{}
	c.companyRegistry = map[string]siteCrawlerFunc{
		"NORDEUS": companies.CrawlNordeus,
		// Add new companies and their functions
	}
	return c
}

func (c *Crawler) CrawlAll() []CompanyCrawlResult {
	/* NOTE: CrawlAll is best-effort.
	It attempts to crawl every registered company even if one crawl fails.
	Each company result contains the successfully crawled job postings
	and the error, if one occurred.
	*/

	var allResults []CompanyCrawlResult

	start := time.Now()

	slog.Info("starting to crawl all company sites")

	for companyName, crawlerFunc := range c.companyRegistry {
		companyStart := time.Now()

		crawledJobs, err := crawlerFunc()

		allResults = append(allResults, CompanyCrawlResult{
			Company:     companyName,
			JobPostings: crawledJobs,
			Err:         err,
		})

		if err != nil {
			slog.Warn("company crawl completed with errors",
				"company", companyName,
				"jobs", len(crawledJobs),
				"duration", time.Since(companyStart),
				"error", err)
			continue
		}

		slog.Debug("company crawl completed", "company", companyName, "jobs", len(crawledJobs), "duration", time.Since(companyStart))
	}

	slog.Info(
		"finished crawling all company sites",
		"companies", len(allResults),
		"duration", time.Since(start),
	)

	return allResults
}

func (c *Crawler) CrawlCompany(company string) (*CompanyCrawlResult, error) {
	crawlerFunc, ok := c.companyRegistry[company]
	if !ok {
		return &CompanyCrawlResult{}, fmt.Errorf("crawler for company %s not found", company)
	}
	jobs, err := crawlerFunc()

	if err != nil {
		slog.Error("error during company crawl", "error", err)
		return nil, err
	}

	return &CompanyCrawlResult{
		Company:     company,
		JobPostings: jobs,
		Err:         err,
	}, nil
}
