package service

import (
	"context"
	"jobfind/model"
	"jobfind/repository"
)

type JobPostingService struct {
	repository repository.JobPostingRepository
}

func NewJobPostingService(repo repository.JobPostingRepository) *JobPostingService {
	return &JobPostingService{
		repository: repo,
	}
}

func (jps *JobPostingService) Sync(ctx context.Context, jobPostings []model.JobPosting) error {
	return nil
}
