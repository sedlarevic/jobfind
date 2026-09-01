package model

type RefreshResult struct {
	Crawled     int32 `json:"crawled"`
	New         int32 `json:"new"`
	Updated     int32 `json:"updated"`
	Deactivated int32 `json:"deactivated"`
}
