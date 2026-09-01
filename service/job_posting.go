package service

import (
	"context"
	"fmt"
	"jobfind/crawler"
	"jobfind/model"
	"jobfind/repository"
	"log/slog"
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

		if !deactivateMissing {
			slog.Warn(
				"skipping missing job deactivation",
				"company", result.Company,
				"reason", "crawl completed with errors",
				"error", result.Err,
			)
		}

		refreshResult, err := jps.repository.RefreshCompany(ctx, result.Company, result.JobPostings, deactivateMissing)

		if err != nil {
			return nil, fmt.Errorf(
				"refresh company %s: %w",
				result.Company,
				err,
			)
		}

		finalResult.Crawled += refreshResult.Crawled
		finalResult.New += refreshResult.New
		finalResult.Updated += refreshResult.Updated
		finalResult.Deactivated += refreshResult.Deactivated

		slog.Debug(
			"company refresh completed",
			"company", result.Company,
			"crawled", refreshResult.Crawled,
			"new", refreshResult.New,
			"updated", refreshResult.Updated,
			"deactivated", refreshResult.Deactivated,
		)
	}

	return finalResult, nil
}

func (jps *JobPostingService) GetAllActive(ctx context.Context) ([]model.JobPosting, error) {
	jobs, err := jps.repository.GetAllActive(ctx)

	if err != nil {
		return nil, fmt.Errorf("get all active job postings: %w", err)
	}

	slog.Debug(
		"active job postings retrieved",
		"jobs", len(jobs),
	)

	return jobs, nil
}
