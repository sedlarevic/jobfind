package companies

import (
	"errors"
	"fmt"
	"jobfind/model"
	"log/slog"
	"strings"

	"github.com/gocolly/colly/v2"
)

func findStringWithWord(strs []string, target string) string {
	for _, str := range strs {
		if strings.Contains(str, target) {
			return str
		}
	}
	return ""
}

func CrawlSwissMarketplaceGroup() ([]model.JobPosting, error) {
	var jobPostings []model.JobPosting
	var siteToVisit = "https://swissmarketplace.group/career/"
	var errs []error

	collector := colly.NewCollector(
		colly.AllowedDomains("www.swissmarketplace.com", "swissmarketplace.com"),
		colly.MaxDepth(2),
	)
	slog.Info("created a collector")
	collector.OnHTML("tr.srJobListEven", func(e *colly.HTMLElement) {
		onClickAttr := e.Attr("onclick")

		splitRes := strings.Split(onClickAttr, "\"")

		if len(splitRes) == 0 {
			return
		}

		absoluteURL := findStringWithWord(splitRes, "https://")

		if absoluteURL == "" {
			return
		}

		slog.Debug("found job link", "url", absoluteURL)

		if err := e.Request.Visit(absoluteURL); err != nil {
			slog.Debug("visit call returned error", "url", absoluteURL, "error", err)
		}
	})

	// on jobposting page event
	collector.OnHTML(".jobad_main_job", func(e *colly.HTMLElement) {
		if e.Request.URL.String() == siteToVisit {
			return
		}

		job := model.JobPosting{
			Title:       strings.TrimSpace(e.ChildText("h1")),
			URL:         e.Request.URL.String(),
			Description: strings.TrimSpace(e.ChildText(".job_sections")),
			CompanyName: "SWISS_MARKETPLACE_GROUP",
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
