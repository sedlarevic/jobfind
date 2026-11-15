package model

import (
	"time"
)

type JobPosting struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	CompanyName  string    `json:"company_name"`
	Description  string    `json:"description"`
	URL          string    `json:"url"`
	FirstSeen    time.Time `json:"first_seen"`
	LastSeen     time.Time `json:"last_seen"`
	IsActive     bool      `json:"is_active"`
	Source       string    `json:"source"`
	Country      string    `json:"country"`
	Category     string    `json:"category"`
	SubCategory  string    `json:"subcategory"`
	WorkingModel string    `json:"working_model"`
	Experience   string    `json:"experience"`
	Location     string    `json:"location,omitempty"`
	Embedding    []float32 `json:"embedding,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

func (jp *JobPosting) String() string {
	return jp.Title + "\n\n" + jp.URL + "\n\n" + jp.Description
}
