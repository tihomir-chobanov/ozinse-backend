package model

// Age_Category represents an age rating category.
// swagger:model Age_Category
type Age_Category struct {
	ID      int    `json:"id" example:"1"`
	Range   string `json:"range" example:"10-12"`
	IconUrl string `json:"icon_url" example:"https://example.com/age-icon.png"`
}
