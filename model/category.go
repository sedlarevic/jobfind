package model

type Category struct {
	Slug          string        `json:"slug"`
	Keywords      []string      `json:"keywords"`
	Subcategories []Subcategory `json:"subcategories"`
}
