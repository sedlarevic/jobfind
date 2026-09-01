package crawler

import (
	"jobfind/crawler/companies"
	"jobfind/model"
	"log"
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

	for companyName, crawlerFunc := range c.companyRegistry {

		crawledJobs, err := crawlerFunc()

		allResults = append(allResults, CompanyCrawlResult{
			Company:     companyName,
			JobPostings: crawledJobs,
			Err:         err,
		})
		log.Printf("crawled job postings of company:\t%v\nnumber of job postings:\t%d\nerrors:%v\n", companyName, len(crawledJobs), err)
	}

	return allResults
}
