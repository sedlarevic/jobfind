package companies

import (
	"errors"
	"fmt"
	"jobfind/model"
	"log/slog"
	"strings"

	"github.com/gocolly/colly/v2"
)

func CrawlNordeus(url string) ([]model.CrawledJob, error) {
	var errs []error
	var crawledJobs []model.CrawledJob

	collector := colly.NewCollector(
		colly.AllowedDomains("www.nordeus.com", "nordeus.com"),
		colly.MaxDepth(2),
	)

	// url event
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)

		if !strings.Contains(absoluteURL, "/open-positions") || strings.Contains(absoluteURL, "/early_talent_program") || absoluteURL == url {
			return
		}

		slog.Debug("found job link", "url", absoluteURL)

		if err := e.Request.Visit(absoluteURL); err != nil {
			slog.Debug("visit call returned error", "url", absoluteURL, "error", err)
		}
	})

	// on jobposting page event
	collector.OnHTML(".main", func(e *colly.HTMLElement) {
		if e.Request.URL.String() == url {
			return
		}

		crawledText := strings.TrimSpace(e.ChildText(".text-content"))
		title := strings.TrimSpace(e.ChildText("h1"))

		if title == "" || crawledText == "" {
			return
		}

		curURL := e.Request.URL.String()
		crawledJob := model.CrawledJob{
			Title:       title,
			Description: crawledText,
			URL:         curURL,
			Country:     "Serbia",
			Location:    "Belgrade",
		}

		crawledJobs = append(crawledJobs, crawledJob)

		slog.Debug("job page crawled", "url", curURL)
	})

	collector.OnRequest(func(r *colly.Request) {
		slog.Debug("visiting page", "url", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		crawlErr := fmt.Errorf("request %s failed: %w", r.Request.URL, err)

		errs = append(errs, crawlErr)

		slog.Warn("request failed", "url", r.Request.URL.String(), "failed", err)
	})

	err := collector.Visit(url)

	if err != nil {
		errs = append(errs, err)
	}

	return crawledJobs, errors.Join(errs...)
}
