package service

import (
	"context"
	"fmt"
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

func (cs *CrawlService) RefreshJobs(ctx context.Context) (*model.RefreshResult, error) {
	perCompanyResult := cs.crawler.CrawlAll()

	refreshResult, err := cs.jobPostingService.Refresh(ctx, perCompanyResult)
	if err != nil {
		return nil, fmt.Errorf("refresh jobs: %w", err)
	}

	return refreshResult, nil
}
