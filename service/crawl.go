package service

import (
	"context"
	"errors"
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

func (cs *CrawlService) RefreshJobs(ctx context.Context) (*model.RefreshResult, error) {
	perCompanyResult := cs.crawler.CrawlAll()

	if len(perCompanyResult) == 0 {
		log.Printf("no company has been crawled! Returning error\n")
		return nil, errors.New("No company site has been crawled.")
	}

	refreshResult, err := cs.jobPostingService.Refresh(ctx, perCompanyResult)

	if err != nil {

		log.Printf("error during Refresh in jps.Refresh\n")
		return nil, err
	}
	log.Printf("returning refreshresult\n")
	return refreshResult, nil
}
