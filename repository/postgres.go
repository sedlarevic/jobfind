package repository

import (
	"context"
	"jobfind/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JobPostingRepository interface {
	GetAllActive(ctx context.Context) ([]model.JobPosting, error)
	Upsert(ctx context.Context, jobPosting model.JobPosting) error
	MarkInactive(ctx context.Context, jobPostingID string) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (pr *PostgresRepository) GetAllActive(ctx context.Context) ([]model.JobPosting, error) {
	return nil, nil
}

func (pr *PostgresRepository) Upsert(ctx context.Context, jobPosting model.JobPosting) error {
	return nil
}

func (pr *PostgresRepository) MarkInactive(ctx context.Context, jobPostingID string) error {
	return nil
}
