package model

// Genre represents a project genre.
// swagger:model Genre
type Genre struct {
	ID      int    `json:"id" example:"1"`
	Name    string `json:"name" example:"Action"`
	IconUrl string `json:"icon_url" example:"https://example.com/icon.png"`
}
