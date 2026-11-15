package resolver

import (
	"encoding/json"
	"fmt"
	"io"
	"jobfind/model"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

const (
	categoryConfigPath     = "config/categories.json"
	experienceConfigPath   = "config/experiences.json"
	workingModelConfigPath = "config/working_models.json"
	sourceConfigPath       = "config/sources.json"
)

type JobDataResolver struct {
	CategoryResolver     *CategoryResolver
	ExperienceResolver   *ExperienceResolver
	WorkingModelResolver *WorkingModelResolver
	SourceResolver       *SourceResolver
}

func NewJobDataResolver() (*JobDataResolver, error) {
	categoryConfig, err := LoadCategoryConfig(categoryConfigPath)
	if err != nil {
		return nil, err
	}

	experienceConfig, err := LoadExperienceConfig(experienceConfigPath)
	if err != nil {
		return nil, err
	}

	workingModelConfig, err := LoadWorkingModelConfig(workingModelConfigPath)
	if err != nil {
		return nil, err
	}

	sourceConfig, err := LoadSourceConfig(sourceConfigPath)
	if err != nil {
		return nil, err
	}

	return &JobDataResolver{
		CategoryResolver:     NewCategoryResolver(categoryConfig),
		ExperienceResolver:   NewExperienceResolver(experienceConfig),
		WorkingModelResolver: NewWorkingModelResolver(workingModelConfig),
		SourceResolver:       NewSourceResolver(sourceConfig),
	}, nil
}

// CATEGORY

type CategoryResolver struct {
	config model.CategoryConfig
}

func NewCategoryResolver(config model.CategoryConfig) *CategoryResolver {
	return &CategoryResolver{config: config}
}

func LoadCategoryConfig(filePath string) (model.CategoryConfig, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return model.CategoryConfig{}, fmt.Errorf("failed to load categories from filepath %s: %w", filePath, err)
	}
	defer jsonFile.Close()

	var categoryConfig model.CategoryConfig

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return model.CategoryConfig{}, fmt.Errorf("failed to read file from filepath %s: %w", filePath, err)
	}

	err = json.Unmarshal(byteValue, &categoryConfig)
	if err != nil {
		return model.CategoryConfig{}, fmt.Errorf("failed to unmarshal bytevalue from filepath %s: %w", filePath, err)
	}

	return categoryConfig, nil
}

func (r *CategoryResolver) Resolve(title string, description string) (string, error) {
	title = strings.ToLower(title)
	description = strings.ToLower(description)

	for _, category := range r.config.Categories {
		match, err := matchesAny(title, category.Keywords)
		if err != nil {
			return "", err
		}

		if match {
			return category.Slug, nil
		}
	}

	for _, category := range r.config.Categories {
		match, err := matchesAny(description, category.Keywords)
		if err != nil {
			return "", err
		}

		if match {
			return category.Slug, nil
		}
	}
	// TODO: Better log message that there is no keyword match
	slog.Debug("No keyword match for", "title", title)

	return r.config.Fallback, nil
}

func (r *CategoryResolver) ResolveSubcategory(categorySlug string, title string, description string) (string, error) {
	title = strings.ToLower(title)
	description = strings.ToLower(description)

	c, ok := r.findCategoryBySlug(categorySlug)

	if !ok {
		return "", fmt.Errorf("category with slug %q not found", categorySlug)
	}

	for _, subcategory := range c.Subcategories {
		match, err := matchesAny(title, subcategory.Keywords)
		if err != nil {
			return "", err
		}

		if match {
			return subcategory.Slug, nil
		}
	}

	for _, subcategory := range c.Subcategories {
		match, err := matchesAny(description, subcategory.Keywords)
		if err != nil {
			return "", err
		}

		if match {
			return subcategory.Slug, nil
		}
	}

	// TODO: Better log message that there is no keyword match
	slog.Debug("No keyword match for", "title", title)

	return r.config.Fallback, nil
}

// EXPERIENCE

type ExperienceResolver struct {
	config model.ExperienceConfig
}

func NewExperienceResolver(config model.ExperienceConfig) *ExperienceResolver {
	return &ExperienceResolver{config: config}
}

func LoadExperienceConfig(filePath string) (model.ExperienceConfig, error) {

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return model.ExperienceConfig{}, fmt.Errorf("failed to load experiences from filepath %s: %w", filePath, err)
	}

	defer jsonFile.Close()

	var config model.ExperienceConfig

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return model.ExperienceConfig{}, fmt.Errorf("failed to read file from filepath %s: %w", filePath, err)
	}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return model.ExperienceConfig{}, fmt.Errorf("failed to unmarshal bytevalue from filepath %s: %w", filePath, err)
	}

	return config, nil
}

func (r *ExperienceResolver) Resolve(title string, description string) (string, error) {
	title = strings.ToLower(title)
	description = strings.ToLower(description)

	for _, experience := range r.config.Experiences {
		match, err := matchesAny(title, experience.Keywords)
		if err != nil {
			return "", err
		}

		if match {
			return experience.Slug, nil
		}
	}

	for _, experience := range r.config.Experiences {
		match, err := matchesAny(description, experience.Keywords)
		if err != nil {
			return "", err
		}

		if match {
			return experience.Slug, nil
		}
	}

	// TODO: Better log message that there is no keyword match
	slog.Debug("No keyword match for", "title", title)

	return r.config.Fallback, nil
}

// WORKING MODEL

type WorkingModelResolver struct {
	config model.WorkingModelConfig
}

func NewWorkingModelResolver(config model.WorkingModelConfig) *WorkingModelResolver {
	return &WorkingModelResolver{config: config}
}

func LoadWorkingModelConfig(filePath string) (model.WorkingModelConfig, error) {

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return model.WorkingModelConfig{}, fmt.Errorf("failed to load working models from filepath %s: %w", filePath, err)
	}

	defer jsonFile.Close()

	var config model.WorkingModelConfig

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return model.WorkingModelConfig{}, fmt.Errorf("failed to read file from filepath %s: %w", filePath, err)
	}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return model.WorkingModelConfig{}, fmt.Errorf("failed to unmarshal bytevalue from filepath %s: %w", filePath, err)
	}

	return config, nil
}

func (r *WorkingModelResolver) Resolve(title string, description string) (string, error) {
	title = strings.ToLower(title)
	description = strings.ToLower(description)

	for _, category := range r.config.WorkingModels {
		match, err := matchesAny(title, category.Keywords)
		if err != nil {
			return "", err
		}

		if match {
			return category.Slug, nil
		}
	}

	for _, category := range r.config.WorkingModels {
		match, err := matchesAny(description, category.Keywords)
		if err != nil {
			return "", err
		}

		if match {
			return category.Slug, nil
		}
	}

	// TODO: Better log message that there is no keyword match
	slog.Debug("No keyword match for", "title", title)

	return r.config.Fallback, nil
}

// SOURCE

type SourceResolver struct {
	config model.SourceConfig
}

func NewSourceResolver(config model.SourceConfig) *SourceResolver {
	return &SourceResolver{config: config}
}

func LoadSourceConfig(filePath string) (model.SourceConfig, error) {

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return model.SourceConfig{}, fmt.Errorf("failed to load sources from filepath %s: %w", filePath, err)
	}

	defer jsonFile.Close()

	var config model.SourceConfig

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return model.SourceConfig{}, fmt.Errorf("failed to read file from filepath %s: %w", filePath, err)
	}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return model.SourceConfig{}, fmt.Errorf("failed to unmarshal bytevalue from filepath %s: %w", filePath, err)
	}

	return config, nil
}

func (r *SourceResolver) Resolve(url string) (string, error) {
	url = strings.ToLower(url)

	for _, source := range r.config.Sources {
		match, err := matchesAny(url, source.Keywords)
		if err != nil {
			return "", err
		}

		if match {
			return source.Slug, nil
		}
	}

	// TODO: Better log message that there is no keyword match
	slog.Debug("No keyword match for", "url", url)

	return r.config.Fallback, nil
}

func matchesAny(text string, patterns []string) (bool, error) {

	for _, p := range patterns {
		matched, err := regexp.MatchString(p, text)
		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}

func (r *CategoryResolver) findCategoryBySlug(slug string) (model.Category, bool) {
	for _, c := range r.config.Categories {
		if c.Slug == slug {
			return c, true
		}
	}

	return model.Category{}, false
}
