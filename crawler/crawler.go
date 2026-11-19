package crawler

import (
	"fmt"
	"jobfind/builder"
	"jobfind/crawler/companies"
	"jobfind/model"
	"jobfind/preprocessing"
	"log/slog"
	"time"
)

type siteCrawlerFunc func(url string) ([]model.CrawledJob, error)

type registeredCompany struct {
	careersURL string
	crawlFunc  siteCrawlerFunc
}

type Crawler struct {
	// NOTE: company name: crawler function
	jobBuilder      *builder.JobBuilder
	companyRegistry map[string]registeredCompany
}

func NewCrawler() (*Crawler, error) {
	jb, err := builder.NewJobBuilder()
	if err != nil {
		return nil, err
	}

	return &Crawler{
		jobBuilder: jb,
		companyRegistry: map[string]registeredCompany{
			"NORDEUS": {
				careersURL: "https://nordeus.com/open-positions/",
				crawlFunc:  companies.CrawlNordeus,
			},
			// Add new companies and their functions
		},
	}, nil
}

func (c *Crawler) CrawlAll() []model.CompanyCrawlResult {
	/* NOTE: CrawlAll makes a best-effort attempt to crawl all pages,
	meaning it will continue running even if individual request fails.
	Each company result contains successfully crawled job postings
	and an error, if one occurred.
	*/

	var allResults []model.CompanyCrawlResult
	var companyJobPostings []model.JobPosting

	start := time.Now()

	slog.Info("starting to crawl all company sites")

	for companyName, regCompany := range c.companyRegistry {
		companyStart := time.Now()

		crawledJobs, err := regCompany.crawlFunc(regCompany.careersURL)
		company := model.Company{Name: companyName, CareersURL: regCompany.careersURL}

		for _, job := range crawledJobs {

			job = preprocessing.NormalizeJob(job)
			jobPosting, err := c.jobBuilder.Build(company, job)
			if err != nil {
				// TODO: HANDLE ERROR
			}

			companyJobPostings = append(companyJobPostings, jobPosting)
		}

		allResults = append(allResults, model.CompanyCrawlResult{
			Company:     company,
			JobPostings: companyJobPostings,
			Err:         err,
		})

		if err != nil {
			slog.Warn("company crawl completed with errors",
				"company", company.Name,
				"jobs", len(companyJobPostings),
				"duration", time.Since(companyStart),
				"error", err)
			continue
		}

		slog.Debug("company crawl completed", "company", company.Name, "jobs", len(companyJobPostings), "duration", time.Since(companyStart))
	}

	slog.Info(
		"finished crawling all company sites",
		"companies", len(allResults),
		"duration", time.Since(start),
	)

	return allResults
}

func (c *Crawler) CrawlCompany(companyName string) (*model.CompanyCrawlResult, error) {
	registeredCompany, ok := c.companyRegistry[companyName]
	if !ok {
		return &model.CompanyCrawlResult{}, fmt.Errorf("crawler for company %s not found", companyName)
	}

	jobs, err := registeredCompany.crawlFunc(registeredCompany.careersURL)
	company := model.Company{Name: companyName, CareersURL: registeredCompany.careersURL}

	var jobPostings []model.JobPosting

	for _, job := range jobs {

		job = preprocessing.NormalizeJob(job)

		jobPosting, err := c.jobBuilder.Build(company, job)
		if err != nil {
			// TODO: HANDLE ERROR
		}

		jobPostings = append(jobPostings, jobPosting)
	}

	if err != nil {
		return nil, err
	}

	return &model.CompanyCrawlResult{
		Company:     company,
		JobPostings: jobPostings,
		Err:         err,
	}, nil
}
