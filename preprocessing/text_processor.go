package preprocessing

import (
	"jobfind/model"
	"strings"
)

func Normalize(s string) string {
	s = strings.ToLower(s)

	return s
}

func NormalizeJob(job model.CrawledJob) model.CrawledJob {
	job.Title = strings.ToLower(job.Title)
	job.Description = strings.ToLower(job.Description)
	job.Country = strings.ToLower(job.Country)
	job.Location = strings.ToLower(job.Location)

	return job
}
