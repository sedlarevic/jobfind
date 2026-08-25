package service

import (
	"context"
	"jobfind/model"
	"jobfind/repository"
)

type JobPostingService struct {
	repository *repository.PostgresRepository
}

func NewJobPostingService(repo *repository.PostgresRepository) *JobPostingService {
	return &JobPostingService{
		repository: repo,
	}
}

func (jps *JobPostingService) Sync(ctx context.Context, jobPostings []model.JobPosting) error {
	return nil
}
