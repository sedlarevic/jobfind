package repository

import (
	"context"
	"fmt"
	"jobfind/model"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JobPostingRepository interface {
	GetAllActive(ctx context.Context) ([]model.JobPosting, error)
	RefreshCompany(ctx context.Context, company string, jobs []model.JobPosting, deactivateMissing bool) (*model.RefreshResult, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (pr *PostgresRepository) RefreshCompany(ctx context.Context, company string, jobs []model.JobPosting, deactivateMissing bool) (*model.RefreshResult, error) {
	res := &model.RefreshResult{}

	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin refresh transaction for %s: %w", company, err)
	}

	defer tx.Rollback(ctx)

	// NOTE: URLs that are crawled
	var urls []string

	for _, job := range jobs {
		urls = append(urls, job.URL)

		var inserted bool
		// NOTE: Upserting a job posting
		err := tx.QueryRow(ctx, `INSERT INTO job_postings (
			company_name,
			title,
			url, 
			description
		)
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (url)
		DO UPDATE SET
		company_name = EXCLUDED.company_name,
		title = EXCLUDED.title,
		description = EXCLUDED.description,
		last_seen = NOW(),
		active = TRUE
		RETURNING xmax = 0;`,
			company, job.Title, job.URL, job.Description).Scan(&inserted)

		if err != nil {
			return nil, fmt.Errorf("upsert job %q for company %s: %w", job.Title, company, err)
		}

		// NOTE: Saving upserting result
		if inserted {
			res.New++
		} else {
			res.Updated++
		}
		res.Crawled++

	}

	// NOTE: Expired job postings deactivation
	if deactivateMissing {
		tag, err := tx.Exec(ctx, `UPDATE job_postings SET active = false
		WHERE company_name = $1 AND active = true AND NOT (url = ANY($2));`, company, urls)

		if err != nil {
			return nil, fmt.Errorf("deactivate missing jobs for company %s: %w", company, err)
		}
		// NOTE: Saving expired job posting deactivation results
		res.Deactivated = int32(tag.RowsAffected())

		slog.Debug(
			"missing jobs deactivated",
			"company", company,
			"deactivated", res.Deactivated,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit refresh transaction for company %s: %w", company, err)
	}

	return res, nil
}

func (pr *PostgresRepository) GetAllActive(ctx context.Context) ([]model.JobPosting, error) {

	rows, err := pr.pool.Query(ctx, "SELECT id, company_name, title, url, description, first_seen, last_seen, active FROM job_postings WHERE active = true;")

	if err != nil {
		return nil, fmt.Errorf("query active job postings: %w", err)
	}

	defer rows.Close()

	var jobs []model.JobPosting

	for rows.Next() {
		var job model.JobPosting

		if err := rows.Scan(
			&job.ID,
			&job.CompanyName,
			&job.Title,
			&job.URL,
			&job.Description,
			&job.FirstSeen,
			&job.LastSeen,
			&job.Active,
		); err != nil {
			return nil, fmt.Errorf("scan active job posting: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate active job postings: %w", err)
	}

	return jobs, nil
}
