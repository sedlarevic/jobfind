package main

import (
	"context"
	"fmt"
	"jobfind/crawler"
	"jobfind/handler"
	"jobfind/repository"
	"jobfind/service"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	// init db
	pool, err := pgxpool.New(
		context.Background(),
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// init repository
	repo := repository.NewPostgresRepository(pool)

	// init services
	jobPostingService := service.NewJobPostingService(repo)
	crawlService := service.NewCrawlService(crawler.NewCrawler(), jobPostingService)

	// init handler
	httpHandler := handler.NewHTTPHandler(jobPostingService, crawlService)

	// setup routes
	mux := httpHandler.SetupRoutes()

	fmt.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
