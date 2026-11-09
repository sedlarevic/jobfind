package model

type CompanyCrawlResult struct {
	Company     Company
	JobPostings []JobPosting
	Err         error
}
