package service

import (
	"context"
	"jobfind/crawler"
	"jobfind/model"
	"jobfind/repository"
	"log"
)

type JobPostingService struct {
	repository repository.JobPostingRepository
}

func NewJobPostingService(repo repository.JobPostingRepository) *JobPostingService {
	return &JobPostingService{
		repository: repo,
	}
}

func (jps *JobPostingService) Refresh(ctx context.Context, perCompanyCrawlResult []crawler.CompanyCrawlResult) (*model.RefreshResult, error) {
	finalResult := &model.RefreshResult{}

	for _, result := range perCompanyCrawlResult {
		deactivateMissing := result.Err == nil

		refreshResult, err := jps.repository.RefreshCompany(ctx, result.Company, result.JobPostings, deactivateMissing)
		if err != nil {
			return nil, err
		}

		finalResult.Crawled += refreshResult.Crawled
		finalResult.New += refreshResult.New
		finalResult.Updated += refreshResult.Updated
		finalResult.Deactivated += refreshResult.Deactivated
		log.Printf("refreshResult of company:%v is %v", result.Company, refreshResult)
	}

	log.Printf("finalResult: %v\n", finalResult)
	return finalResult, nil
}

func (jps *JobPostingService) GetAllActive(ctx context.Context) ([]model.JobPosting, error) {
	jobs, err := jps.repository.GetAllActive(ctx)

	if err != nil {
		return nil, err
	}

	return jobs, nil
}
