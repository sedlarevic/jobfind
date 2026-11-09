package model

import "time"

type DBJobPosting struct {
	ID             int64
	Title          string
	CompanyID      int64
	Description    string
	URL            string
	FirstSeen      time.Time
	LastSeen       time.Time
	IsActive       bool
	Source         string
	Location       *string
	CategoryID     int32
	CountryID      int32
	WorkingModelID int32
	ExperienceID   int32
	Embedding      []float32
	CreatedAt      time.Time
}

func (dbjp *DBJobPosting) String() string {
	return dbjp.Title + "\n\n" + dbjp.URL + "\n\n" + dbjp.Description
}
