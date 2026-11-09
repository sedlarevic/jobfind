package repository

import (
	"context"
	"fmt"
	"jobfind/model"
)

type subcategoryKey struct {
	CategoryID int32
	Name       string
}

type jobReferenceIDs struct {
	subcategoryID  int32
	experienceID   int32
	sourceID       int32
	countryID      int32
	workingModelID int32
}

func (pr *PostgresRepository) resolveJobReferenceIDs(job model.JobPosting) (jobReferenceIDs, error) {

	categoryID, ok := pr.categories[string(job.Category)]
	if !ok {
		return jobReferenceIDs{}, fmt.Errorf("unknown category: %s", job.Category)
	}

	key := subcategoryKey{
		CategoryID: categoryID,
		Name:       string(job.SubCategory),
	}

	subcategoryID, ok := pr.subcategories[key]
	if !ok {
		return jobReferenceIDs{}, fmt.Errorf(
			"unknown subcategory %s for category %s",
			job.SubCategory,
			job.Category,
		)
	}

	sourceID, ok := pr.sources[string(job.Source)]
	if !ok {
		return jobReferenceIDs{}, fmt.Errorf("unknown source: %s", job.Source)
	}

	countryID, ok := pr.countries[string(job.Country)]
	if !ok {
		return jobReferenceIDs{}, fmt.Errorf("unknown country: %s", job.Country)
	}

	workingModelID, ok := pr.workingModels[string(job.WorkingModel)]
	if !ok {
		return jobReferenceIDs{}, fmt.Errorf("unknown working model: %s", job.WorkingModel)
	}

	experienceID, ok := pr.experiences[string(job.Experience)]
	if !ok {
		return jobReferenceIDs{}, fmt.Errorf("unknown experience: %s", job.Experience)
	}

	return jobReferenceIDs{
		subcategoryID:  subcategoryID,
		sourceID:       sourceID,
		countryID:      countryID,
		workingModelID: workingModelID,
		experienceID:   experienceID,
	}, nil

}

func (pr *PostgresRepository) loadLookups(ctx context.Context) error {
	if err := pr.loadCountriesLookup(ctx); err != nil {
		return err
	}

	if err := pr.loadSourcesLookup(ctx); err != nil {
		return err
	}

	if err := pr.loadWorkingModelsLookup(ctx); err != nil {
		return err
	}

	if err := pr.loadExperiencesLookup(ctx); err != nil {
		return err
	}

	if err := pr.loadCategoriesLookup(ctx); err != nil {
		return err
	}

	if err := pr.loadSubcategoriesLookup(ctx); err != nil {
		return err
	}

	return nil
}

func (pr *PostgresRepository) loadCountriesLookup(ctx context.Context) error {
	rows, err := pr.pool.Query(ctx, `select id, name from countries`)

	if err != nil {
		return fmt.Errorf("countries table load failed: %w", err)
	}

	for rows.Next() {
		var name string
		var id int32

		err = rows.Scan(&id, &name)
		if err != nil {
			return fmt.Errorf("countries row read failed: %w", err)
		}

		pr.countries[name] = id
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("countries rows iteration failed: %w", err)
	}

	rows.Close()

	return nil
}

func (pr *PostgresRepository) loadSourcesLookup(ctx context.Context) error {
	rows, err := pr.pool.Query(ctx, `select id, name from sources`)

	if err != nil {
		return fmt.Errorf("sources table load failed: %w", err)
	}

	for rows.Next() {
		var name string
		var id int32

		err = rows.Scan(&id, &name)
		if err != nil {
			return fmt.Errorf("sources row read failed: %w", err)
		}

		pr.sources[name] = id
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("sources rows iteration failed: %w", err)
	}

	rows.Close()

	return nil
}

func (pr *PostgresRepository) loadWorkingModelsLookup(ctx context.Context) error {
	rows, err := pr.pool.Query(ctx, `SELECT id, name FROM working_models`)

	if err != nil {
		return fmt.Errorf("working models table load failed: %w", err)
	}

	for rows.Next() {
		var name string
		var id int32

		err = rows.Scan(&id, &name)

		if err != nil {
			return fmt.Errorf("working models row read failed: %w", err)
		}

		pr.workingModels[name] = id
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("working models rows iteration failed: %w", err)
	}

	rows.Close()

	return nil
}

func (pr *PostgresRepository) loadExperiencesLookup(ctx context.Context) error {
	rows, err := pr.pool.Query(ctx, `SELECT id, name FROM experiences`)

	if err != nil {
		return fmt.Errorf("experience table load failed: %w", err)
	}

	for rows.Next() {
		var name string
		var id int32

		err = rows.Scan(&id, &name)
		if err != nil {
			return fmt.Errorf("experience row read failed: %w", err)
		}

		pr.experiences[name] = id
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("experience rows iteration failed: %w", err)
	}

	rows.Close()

	return nil
}

func (pr *PostgresRepository) loadCategoriesLookup(ctx context.Context) error {
	rows, err := pr.pool.Query(ctx, `SELECT id, name FROM job_categories`)

	if err != nil {
		return fmt.Errorf("category table load failed: %w", err)
	}

	for rows.Next() {
		var name string
		var id int32

		err = rows.Scan(&id, &name)
		if err != nil {
			return fmt.Errorf("category row read failed: %w", err)
		}

		pr.categories[name] = id
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("category rows iteration failed: %w", err)
	}

	rows.Close()

	return nil
}

func (pr *PostgresRepository) loadSubcategoriesLookup(ctx context.Context) error {
	rows, err := pr.pool.Query(ctx, `SELECT id, name, category_id FROM job_subcategories`)

	if err != nil {
		return fmt.Errorf("subcategory table load failed: %w", err)
	}

	for rows.Next() {
		var subcategoryKey subcategoryKey
		var id int32

		err = rows.Scan(&id, &subcategoryKey.Name, &subcategoryKey.CategoryID)
		if err != nil {
			return fmt.Errorf("subcategory row read failed: %w", err)
		}

		pr.subcategories[subcategoryKey] = id
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("subcategory rows iteration failed: %w", err)
	}

	rows.Close()

	return nil
}
