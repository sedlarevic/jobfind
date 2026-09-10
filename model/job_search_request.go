package model

type JobSearchRequest struct {
	Embedding []float32 `json:"embedding"`
	Limit     int32       `json:"limit"`
}
