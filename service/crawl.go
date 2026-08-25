package service

import (
	"context"
	"jobfind/crawler"
	"jobfind/model"
	"log"
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

func (cs *CrawlService) CrawlAll() ([]model.JobPosting, error) {
	crawledJobs, err := cs.crawler.CrawlAll()

	if err != nil {
		log.Printf("Error occured during crawling: %v", err)
		return nil, err
	}

	err = cs.jobPostingService.Sync(context.Background(), crawledJobs)

	if err != nil {
		log.Printf("Error occured during job posting synchronizing: %v", err)
		return nil, err
	}

	return crawledJobs, nil
}
