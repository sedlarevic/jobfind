package main

import (
	"context"
	"fmt"
	"jobfind/crawler"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/jackc/pgx/v5"
)

func main() {

	// urlExample := "postgres://username:password@localhost:5432/database_name"

	//starting db
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	//starting crawler
	router := mux.NewRouter().StrictSlash(true)
	crawler := crawler.NewCrawler()

	router.HandleFunc("/crawl", crawler.CrawlAllSites).Methods("GET")

	fmt.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", router))
}
