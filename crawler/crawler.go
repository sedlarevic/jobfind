package crawler

import (
	"net/http"

	"github.com/gocolly/colly/v2"
)

func Crawl(w http.ResponseWriter, r http.Response) {

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(100),
	)
	urls := []string{
		"https://jobs.eu.lever.co/symphony",
		"https://swissmarketplace.group/career/belgrade/",
		"https://nordeus.com/open-positions/"
	}

}
