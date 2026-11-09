package repository

import (
	"context"
	"fmt"
	"jobfind/model"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type JobPostingRepository interface {
	GetAllActive(ctx context.Context) ([]model.JobPosting, error)
	RefreshCompany(
		ctx context.Context,
		company model.Company,
		jobs []model.JobPosting,
		deactivateMissing bool) (*model.RefreshResult, error)
	SearchJobs(
		ctx context.Context,
		embedding []float32,
		limit int32) ([]model.JobPosting, error)
	GetJobByID(ctx context.Context, id int64) (*model.JobPosting, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool

	experiences   map[string]int32
	workingModels map[string]int32
	sources       map[string]int32
	countries     map[string]int32
	categories    map[string]int32
	subcategories map[subcategoryKey]int32
}

func NewPostgresRepository(ctx context.Context, pool *pgxpool.Pool) (*PostgresRepository, error) {
	pr := &PostgresRepository{
		pool: pool,

		experiences:   make(map[string]int32),
		workingModels: make(map[string]int32),
		sources:       make(map[string]int32),
		countries:     make(map[string]int32),
		categories:    make(map[string]int32),
		subcategories: make(map[subcategoryKey]int32),
	}

	if err := pr.loadLookups(ctx); err != nil {
		return nil, fmt.Errorf("loading of lookups failed during initialization of postgres repository struct: %w", err)
	}

	return pr, nil
}

func (pr *PostgresRepository) RefreshCompany(ctx context.Context, company model.Company, jobs []model.JobPosting, deactivateMissing bool) (*model.RefreshResult, error) {
	res := &model.RefreshResult{}

	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin refresh transaction for %s: %w", company, err)
	}

	defer tx.Rollback(ctx)

	// NOTE: Upsert company. If exists, return ID, if not create a row and return ID.

	var companyID int64

	err = tx.QueryRow(ctx, `INSERT INTO companies (
		name,
		careers_url
	) VALUES ($1, $2) 
	ON CONFLICT (name) 
	DO UPDATE SET careers_url = EXCLUDED.careers_url
	RETURNING id;`, company.Name, company.CareersURL).Scan(&companyID)

	if err != nil {
		return nil, fmt.Errorf("company upsert failed for company: %s: %w", company.Name, err)
	}

	var urls []string

	for _, job := range jobs {
		// NOTE: URLs that are crawled
		urls = append(urls, job.URL)

		refs, err := pr.resolveJobReferenceIDs(job)

		if err != nil {
			return nil, fmt.Errorf("error during job reference ids resolution for job %s at company %s: %w", job.Title, job.CompanyName, err)
		}

		var inserted bool
		// NOTE: Upserting a job posting
		err = tx.QueryRow(ctx, `INSERT INTO job_postings (
			title,
			url,
			description,
			first_seen,
			last_seen,
			is_active,
			source_id,
			location,
			subcategory_id,
			country_id,
			working_model_id,
			experience_id,
			company_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) 
		ON CONFLICT (url)
		DO UPDATE SET
		title = EXCLUDED.title,
		description = EXCLUDED.description,
		last_seen = NOW(),
		is_active = TRUE,
		source_id = EXCLUDED.source_id,
		location = EXCLUDED.location,
		subcategory_id = EXCLUDED.subcategory_id,
		country_id = EXCLUDED.country_id,
		working_model_id = EXCLUDED.working_model_id,
		experience_id = EXCLUDED.experience_id,
		company_id = EXCLUDED.company_id,
		embedding = CASE
			WHEN job_postings.title IS DISTINCT FROM EXCLUDED.title
				OR job_postings.description IS DISTINCT FROM EXCLUDED.description
			THEN NULL
			ELSE job_postings.embedding
		END
		RETURNING xmax = 0;`,
			job.Title, job.URL, job.Description, job.FirstSeen,
			job.LastSeen, job.IsActive, refs.sourceID,
			job.Location, refs.subcategoryID, refs.countryID,
			refs.workingModelID, refs.experienceID, companyID).Scan(&inserted)

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
		tag, err := tx.Exec(ctx, `UPDATE job_postings SET is_active = false
		WHERE company_id = $1 AND is_active = true AND NOT (url = ANY($2));`, companyID, urls)

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
			&job.IsActive,
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

func (pr *PostgresRepository) SearchJobs(ctx context.Context, embedding []float32, limit int32) ([]model.JobPosting, error) {
	queryVector := pgvector.NewVector(embedding)

	rows, err := pr.pool.Query(ctx, `
	SELECT 
		id,
		company_name,
		title,
		url,
		description,
		first_seen,
		last_seen,
		active
		
		FROM job_postings
		
		WHERE active = true AND embedding IS NOT NULL
		ORDER BY embedding <=> $1 LIMIT $2;`, queryVector, limit)
	if err != nil {
		return nil, fmt.Errorf("query search job postings: %w", err)
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
			&job.IsActive,
		); err != nil {
			return nil, fmt.Errorf("scan searched job postings: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate searched job postings: %w", err)
	}

	return jobs, nil
}

func (pr *PostgresRepository) GetJobByID(ctx context.Context, id int64) (*model.JobPosting, error) {

	var job model.JobPosting

	err := pr.pool.QueryRow(ctx, "SELECT id, company_name, title, url, description, first_seen, last_seen, active FROM job_postings WHERE id = $1;", id).Scan(
		&job.ID,
		&job.CompanyName,
		&job.Title,
		&job.URL,
		&job.Description,
		&job.FirstSeen,
		&job.LastSeen,
		&job.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("get job posting by id: %w", err)
	}

	return &job, nil
}
