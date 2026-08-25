package main

import (
	"context"
	"fmt"
	"jobfind/crawler"
	"jobfind/handler"
	"jobfind/repository"
	"jobfind/service"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	//starting db
	pool, err := pgxpool.New(
		context.Background(),
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := repository.NewPostgresRepository(pool)

	jobPostingService := service.NewJobPostingService(repo)
	crawlService := service.NewCrawlService(crawler.NewCrawler(), jobPostingService)

	httpHandler := handler.NewHTTPHandler(jobPostingService, crawlService)

	mux := httpHandler.SetupRoutes()

	fmt.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
