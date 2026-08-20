package model

import (
	"time"
)

type JobPosting struct {
	Id          int64     `json:"id"`
	CompanyName string    `json:"companyName"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	FirstSeen   time.Time `json:"firstSeen"`
	LastSeen    time.Time `json:"LastSeen"`
	Active      bool      `json:"active"`
}

func (jp *JobPosting) String() string {
	return jp.Title + "\n\n" + jp.URL + "\n\n" + jp.Description
}
