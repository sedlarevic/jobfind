package model

import (
	"time"
)

type JobPosting struct {
	ID           int64        `json:"id"`
	Title        string       `json:"title"`
	CompanyName  string       `json:"companyName"`
	Description  string       `json:"description"`
	URL          string       `json:"url"`
	FirstSeen    time.Time    `json:"firstSeen"`
	LastSeen     time.Time    `json:"lastSeen"`
	IsActive     bool         `json:"isActive"`
	Source       Source       `json:"source"`
	Country      string       `json:"country"`
	Category     Category     `json:"category"`
	SubCategory  SubCategory  `json:"subCategory"`
	WorkingModel WorkingModel `json:"workingModel"`
	Experience   Experience   `json:"experience"`
	Location     string       `json:"location,omitempty"`
	Embedding    []float32    `json:"embedding,omitempty"`
	CreatedAt    time.Time    `json:"createdAt"`
}

func (jp *JobPosting) String() string {
	return jp.Title + "\n\n" + jp.URL + "\n\n" + jp.Description
}
