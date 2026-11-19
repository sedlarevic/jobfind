package builder

import (
	"jobfind/model"
	"jobfind/resolver"
	"time"
)

type JobBuilder struct {
	resolver *resolver.JobDataResolver
}

func NewJobBuilder() (*JobBuilder, error) {
	r, err := resolver.NewJobDataResolver()
	if err != nil {
		return nil, err
	}

	return &JobBuilder{
		resolver: r,
	}, nil
}

func (jb *JobBuilder) Build(company model.Company, crawledJob model.CrawledJob) (model.JobPosting, error) {
	// TODO: Handle errors!
	category, err := jb.resolver.CategoryResolver.Resolve(crawledJob.Title, crawledJob.Description)
	if err != nil {
		return model.JobPosting{}, err
	}

	subcategory, err := jb.resolver.CategoryResolver.ResolveSubcategory(category, crawledJob.Title, crawledJob.Description)
	if err != nil {

		return model.JobPosting{}, err
	}

	workingModel, err := jb.resolver.WorkingModelResolver.Resolve(crawledJob.Title, crawledJob.Description)
	if err != nil {

		return model.JobPosting{}, err
	}

	experience, err := jb.resolver.ExperienceResolver.Resolve(crawledJob.Title, crawledJob.Description)
	if err != nil {
		return model.JobPosting{}, err
	}

	source, err := jb.resolver.SourceResolver.Resolve(crawledJob.URL)
	if err != nil {
		return model.JobPosting{}, err
	}

	now := time.Now()

	return model.JobPosting{
		Title:        crawledJob.Title,
		CompanyName:  company.Name,
		Description:  crawledJob.Description,
		URL:          crawledJob.URL,
		FirstSeen:    now,
		LastSeen:     now,
		IsActive:     true,
		Source:       source,
		Country:      crawledJob.Country,
		Location:     crawledJob.Location,
		Category:     category,
		SubCategory:  subcategory,
		WorkingModel: workingModel,
		Experience:   experience,
	}, nil
}
