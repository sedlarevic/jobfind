package repository

import (
	"context"
	"jobfind/model"
	"log"

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
	log.Printf("in Repository.RefreshCompany for company: %v\n", company)
	res := &model.RefreshResult{}

	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		log.Printf("couldnt create transaction\n")
		return nil, err
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

			log.Printf("error after executing upsert query, %v", err)
			return nil, err
		}

		// NOTE: Saving upserting result
		if inserted {
			res.New++
		} else {
			res.Updated++
		}
		res.Crawled++
		log.Printf("finished upserting a job, inserted:%v\n", inserted)
	}

	// NOTE: Expired job postings deactivation
	if deactivateMissing {
		tag, err := tx.Exec(ctx, `UPDATE job_postings SET active = false
		WHERE company_name = $1 AND active = true AND NOT (url = ANY($2));`, company, urls)

		if err != nil {
			return nil, err
		}
		// NOTE: Saving expired job posting deactivation results
		res.Deactivated = int32(tag.RowsAffected())

		log.Printf("deactivating expired jobs, deactivated:%v\n", res.Deactivated)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("couldnt commit for some reason, %v\n", err)
		return nil, err
	}
	log.Printf("successful refresh, result: %v\n", res)
	return res, nil
}

func (pr *PostgresRepository) GetAllActive(ctx context.Context) ([]model.JobPosting, error) {

	rows, err := pr.pool.Query(ctx, "SELECT * FROM job_postings WHERE active = true;")

	if err != nil {
		return nil, err
	}

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
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
