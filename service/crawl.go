package service

import (
	"context"
	"fmt"
	"jobfind/crawler"
	"jobfind/model"
	"strings"
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

func (cs *CrawlService) CrawlCompany(ctx context.Context, company string) (*crawler.CompanyCrawlResult, error) {
	company = strings.ToUpper(company)

	crawlResult, err := cs.crawler.CrawlCompany(company)

	if err != nil {
		return nil, fmt.Errorf("crawl company: %w", err)
	}

	return crawlResult, nil
}
