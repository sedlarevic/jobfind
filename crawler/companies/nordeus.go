package companies

import (
	"errors"
	"fmt"
	"jobfind/model"
	"jobfind/resolver"
	"log/slog"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func CrawlNordeus() ([]model.JobPosting, error) {
	var jobPostings []model.JobPosting
	var siteToVisit = "https://nordeus.com/open-positions/"
	var errs []error

	collector := colly.NewCollector(
		colly.AllowedDomains("www.nordeus.com", "nordeus.com"),
		colly.MaxDepth(2),
	)

	// url event
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)

		if !strings.Contains(absoluteURL, "/open-positions") || strings.Contains(absoluteURL, "/early_talent_program") || absoluteURL == siteToVisit {
			return
		}

		slog.Debug("found job link", "url", absoluteURL)

		if err := e.Request.Visit(absoluteURL); err != nil {
			slog.Debug("visit call returned error", "url", absoluteURL, "error", err)
		}
	})

	// on jobposting page event
	collector.OnHTML(".main", func(e *colly.HTMLElement) {
		if e.Request.URL.String() == siteToVisit {
			return
		}
		crawledText := strings.TrimSpace(e.ChildText(".text-content"))
		title := strings.TrimSpace(e.ChildText("h1"))

		resolver, err := resolver.NewJobDataResolver()
		if err != nil {
			// TODO: DO SOMETHING
		}
		// TODO: Handle errors!
		category, err := resolver.CategoryResolver.Resolve(title, crawledText)
		subCategory, err := resolver.CategoryResolver.ResolveSubcategory(category, title, crawledText)
		workingModel, err := resolver.WorkingModelResolver.Resolve(title, crawledText)
		experience, err := resolver.ExperienceResolver.Resolve(title, crawledText)
		source, err := resolver.SourceResolver.Resolve(siteToVisit)

		job := model.JobPosting{
			Title:        title,
			CompanyName:  "NORDEUS",
			Description:  crawledText,
			URL:          e.Request.URL.String(),
			FirstSeen:    time.Now(),
			LastSeen:     time.Now(),
			IsActive:     true,
			Source:       source,
			Country:      "Serbia",
			Category:     category,
			SubCategory:  subCategory,
			WorkingModel: workingModel,
			Experience:   experience,
			Location:     "Belgrade, Serbia",
		}

		jobPostings = append(jobPostings, job)

		slog.Debug("job page crawled", "url", job.URL)
	})

	collector.OnRequest(func(r *colly.Request) {
		slog.Debug("visiting page", "url", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		crawlErr := fmt.Errorf("request %s failed: %w", r.Request.URL, err)

		errs = append(errs, crawlErr)

		slog.Warn("request failed", "url", r.Request.URL.String(), "failed", err)
	})

	err := collector.Visit(siteToVisit)

	if err != nil {
		errs = append(errs, err)
	}

	return jobPostings, errors.Join(errs...)
}
