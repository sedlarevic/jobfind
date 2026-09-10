package main

import (
	"context"
	"jobfind/crawler"
	"jobfind/cv"
	"jobfind/handler"
	"jobfind/repository"
	"jobfind/service"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvector "github.com/pgvector/pgvector-go/pgx"
)

func main() {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	ctx := context.Background()

	// init db
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("database config failed", "error", err)
		os.Exit(1)
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxvector.RegisterTypes(ctx, conn)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("database initialization failed", "error", err)
		os.Exit(1)
	}
	if err != nil {
		slog.Error("database initialization failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}

	slog.Info("database connection established")

	// init repository
	repo := repository.NewPostgresRepository(pool)

	// init services
	jobPostingService := service.NewJobPostingService(repo)
	crawlService := service.NewCrawlService(crawler.NewCrawler(), jobPostingService)
	cvService := service.NewCVService(cv.NewCVReader())

	// init handler
	httpHandler := handler.NewHTTPHandler(jobPostingService, crawlService, cvService)

	// setup routes
	mux := httpHandler.SetupRoutes()

	slog.Info("server starting", "port", 8081)

	if err := http.ListenAndServe(":8081", mux); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
