package service

import (
	"context"
	"jobfind/crawler"
	"jobfind/model"
)

type CrawlService struct {
	crawler           *crawler.Crawler
	jobPostingService *JobPostingService
}

func NewCrawlService(crawler *crawler.Crawler, jps *JobPostingService) *CrawlService {
	return &CrawlService{
		crawler:           crawler,
		jobPostingService: jps,
	}
}

func (cs *CrawlService) CrawlAll(ctx context.Context) ([]model.JobPosting, error) {
	crawledJobs, crawlErr := cs.crawler.CrawlAll()

	syncErr := cs.jobPostingService.Sync(ctx, crawledJobs)

	if syncErr != nil {
		return nil, syncErr
	}

	return crawledJobs, crawlErr
}
