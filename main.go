package main

import (
	"fmt"
	"jobfind/crawler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	crawler := crawler.NewCrawler()

	router.HandleFunc("/crawl", crawler.CrawlAllSites).Methods("GET")

	fmt.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", router))
}
