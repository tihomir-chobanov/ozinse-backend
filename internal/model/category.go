package model

// Category represents a content category used to tag projects.
// swagger:model Category
type Category struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Drama"`
}
