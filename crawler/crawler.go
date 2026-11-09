package crawler

import (
	"fmt"
	"jobfind/crawler/companies"
	"jobfind/model"
	"log/slog"
	"time"
)

type siteCrawlerFunc func() ([]model.JobPosting, error)

type Crawler struct {
	// NOTE: company name: crawler function
	companyRegistry map[string]registeredCompany
}

func NewCrawler() *Crawler {
	return &Crawler{
		companyRegistry: map[string]registeredCompany{
			"NORDEUS": {
				careersURL: "https://nordeus.com/open-positions/",
				crawlFunc:  companies.CrawlNordeus,
			},
			// Add new companies and their functions
		},
	}
}

type registeredCompany struct {
	careersURL string
	crawlFunc  siteCrawlerFunc
}

func (c *Crawler) CrawlAll() []model.CompanyCrawlResult {
	/* NOTE: CrawlAll is best-effort.
	It attempts to crawl every registered company even if one crawl fails.
	Each company result contains the successfully crawled job postings
	and the error, if one occurred.
	*/

	var allResults []model.CompanyCrawlResult

	start := time.Now()

	slog.Info("starting to crawl all company sites")

	for companyName, company := range c.companyRegistry {
		companyStart := time.Now()

		crawledJobs, err := company.crawlFunc()

		allResults = append(allResults, model.CompanyCrawlResult{
			Company: model.Company{
				Name:       companyName,
				CareersURL: company.careersURL,
			},
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

func (c *Crawler) CrawlCompany(companyName string) (*model.CompanyCrawlResult, error) {
	company, ok := c.companyRegistry[companyName]
	if !ok {
		return &model.CompanyCrawlResult{}, fmt.Errorf("crawler for company %s not found", company)
	}

	jobs, err := company.crawlFunc()

	if err != nil {
		return nil, err
	}

	return &model.CompanyCrawlResult{
		Company:     model.Company{Name: companyName, CareersURL: company.careersURL},
		JobPostings: jobs,
		Err:         err,
	}, nil
}
